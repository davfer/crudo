package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/davfer/crudo"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/inmemory"
	"github.com/davfer/crudo/notifier"
	"github.com/davfer/go-specification"
)

const (
	Loaded   = "loaded"
	Unloaded = "unloaded"
	Added    = "added"
	Updated  = "updated"
	Deleted  = "deleted"
)

type HydrateFunc[K entity.Entity] func(ctx context.Context, entity K) (K, error)
type RefreshPolicy string

const (
	RefreshPolicyNone         RefreshPolicy = "none"
	RefreshPolicyReadAll      RefreshPolicy = "read-all"
	RefreshPolicyWriteAll     RefreshPolicy = "write-all"
	RefreshPolicyReadWriteAll RefreshPolicy = "read-write-all"
)

type ProxyStore[K entity.Entity] struct {
	remoteRepository crudo.Repository[K]
	localRepository  crudo.Repository[K]
	RefreshPolicy    RefreshPolicy
	notifier         *notifier.TopicCallbackNotifier[K]
	Hydrate          HydrateFunc[K]
}

func NewProxyStore[K entity.Entity]() *ProxyStore[K] {
	return &ProxyStore[K]{
		RefreshPolicy: RefreshPolicyNone,
		notifier:      notifier.NewTopicCallbackNotifier[K]([]string{Added, Updated, Deleted, Loaded, Unloaded}),
	}
}

func (r *ProxyStore[K]) On(event string, observer notifier.ObserverCallback[K]) error {
	return r.notifier.Attach(event, observer)
}

func (r *ProxyStore[K]) OnHydrate(onHydrate HydrateFunc[K]) {
	r.Hydrate = onHydrate
}

func (r *ProxyStore[K]) Start(ctx context.Context, onBootstrap func(ctx context.Context) error) error {
	if r.remoteRepository == nil {
		return fmt.Errorf("store not loaded")
	}

	return r.remoteRepository.Start(ctx, onBootstrap)
}

func (r *ProxyStore[K]) Create(ctx context.Context, e K) (K, error) {
	if r.remoteRepository == nil {
		return e, fmt.Errorf("store not loaded")
	}

	if !e.GetId().IsEmpty() {
		return e, fmt.Errorf("entity already with id")
	}

	e, err := r.remoteRepository.Create(ctx, e)
	if err != nil {
		return e, fmt.Errorf("could not insert entity: %w", err)
	}
	// hydrate persisted entity
	if r.Hydrate != nil {
		e, err = r.Hydrate(ctx, e)
		if err != nil {
			return e, fmt.Errorf("could not hydrate entity: %w", err)
		}
	}
	if _, err = r.localRepository.Create(ctx, e); err != nil {
		return e, fmt.Errorf("could not insert entity locally: %w", err)
	}
	if err = r.notifier.Notify(ctx, Added, e); err != nil {
		return e, fmt.Errorf("could not notify entity add: %w", err)
	}

	return e, nil
}

func (r *ProxyStore[K]) Read(ctx context.Context, id entity.Id) (e K, err error) {
	if r.remoteRepository == nil {
		return e, fmt.Errorf("store not loaded")
	}

	e, err = r.localRepository.Read(ctx, id)
	if err != nil && !errors.Is(err, entity.ErrEntityNotFound) {
		err = fmt.Errorf("could not read entity: %w", err)
		return
	} else if err == nil {
		return
	}

	e, err = r.remoteRepository.Read(ctx, id)
	if err != nil && !errors.Is(err, entity.ErrEntityNotFound) {
		err = fmt.Errorf("could not read entity: %w", err)
	} else if err == nil {
		r.localRepository.Create(ctx, e)
		if err = r.notifier.Notify(ctx, Loaded, e); err != nil {
			return
		}
	}

	return
}

func (r *ProxyStore[K]) Match(ctx context.Context, c specification.Criteria) ([]K, error) {
	if r.remoteRepository == nil {
		return []K{}, fmt.Errorf("store not loaded")
	}

	return r.localRepository.Match(ctx, c)
}

func (r *ProxyStore[K]) MatchOne(ctx context.Context, c specification.Criteria) (K, error) {
	if r.remoteRepository == nil {
		return *new(K), fmt.Errorf("store not loaded")
	}

	return r.localRepository.MatchOne(ctx, c)
}

func (r *ProxyStore[K]) ReadAll(ctx context.Context) ([]K, error) {
	if r.remoteRepository == nil {
		return []K{}, fmt.Errorf("store not loaded")
	}

	return r.localRepository.ReadAll(ctx)
}

func (r *ProxyStore[K]) Update(ctx context.Context, entity K) error {
	if r.remoteRepository == nil {
		return fmt.Errorf("store not loaded")
	}

	if err := r.remoteRepository.Update(ctx, entity); err != nil {
		return fmt.Errorf("could not update entity remotelly: %w", err)
	}
	if err := r.localRepository.Update(ctx, entity); err != nil {
		return fmt.Errorf("could not update entity locally: %w", err)
	}
	if err := r.notifier.Notify(ctx, Updated, entity); err != nil {
		return fmt.Errorf("could not notify entity update: %w", err)
	}

	return nil
}

func (r *ProxyStore[K]) Delete(ctx context.Context, entity K) error {
	if r.remoteRepository == nil {
		return fmt.Errorf("store not loaded")
	}

	if err := r.remoteRepository.Delete(ctx, entity); err != nil {
		return fmt.Errorf("could not delete entity remotelly: %w", err)
	}
	if err := r.localRepository.Delete(ctx, entity); err != nil {
		return fmt.Errorf("could not delete entity locally: %w", err)
	}
	if err := r.notifier.Notify(ctx, Deleted, entity); err != nil {
		return fmt.Errorf("could not notify entity delete: %w", err)
	}

	return nil
}

func (r *ProxyStore[K]) Load(ctx context.Context, repo crudo.Repository[K]) error {
	r.remoteRepository = repo

	if r.localRepository != nil {
		return fmt.Errorf("entities already loaded")
	}

	entities, err := r.remoteRepository.ReadAll(ctx)
	if err != nil {
		return fmt.Errorf("could not load Entities: %w", err)
	}

	r.localRepository = inmemory.NewRepository(entities)
	for _, d := range entities {
		if r.Hydrate != nil {
			d, err = r.Hydrate(ctx, d)
			if err != nil {
				return fmt.Errorf("could not hydrate entity: %w", err)
			}
		}

		r.notifier.Notify(ctx, Loaded, d)
	}

	return nil
}
func (r *ProxyStore[K]) Refresh(ctx context.Context) error {
	if r.remoteRepository == nil {
		return fmt.Errorf("store not loaded")
	}

	if r.RefreshPolicy == RefreshPolicyNone {
		return nil
	}

	remoteEntities, err := r.remoteRepository.ReadAll(ctx)
	if err != nil {
		return fmt.Errorf("could not load remote Entities: %w", err)
	}

	localEntities, err := r.localRepository.ReadAll(ctx)
	if err != nil {
		return fmt.Errorf("could not load local Entities: %w", err)
	}

	if r.RefreshPolicy == RefreshPolicyReadAll || r.RefreshPolicy == RefreshPolicyReadWriteAll {
		for _, d := range remoteEntities {
			if !entity.Contains(localEntities, d) {
				if _, err = r.localRepository.Create(ctx, d); err != nil {
					return err
				}
				if err = r.notifier.Notify(ctx, Loaded, d); err != nil {
					return err
				}
			}
		}
	}

	if r.RefreshPolicy == RefreshPolicyWriteAll || r.RefreshPolicy == RefreshPolicyReadWriteAll {
		for _, d := range localEntities {
			if !entity.Contains(remoteEntities, d) {
				_, err = r.remoteRepository.Create(ctx, d)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// unload not found
		for _, d := range localEntities {
			if !entity.Contains(remoteEntities, d) {
				if err = r.localRepository.Delete(ctx, d); err != nil {
					return err
				}
				if err = r.notifier.Notify(ctx, Unloaded, d); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
