package store

import (
	"context"
	"github.com/davfer/crudo"
	"github.com/davfer/crudo/criteria"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/inmemory"
	"github.com/davfer/crudo/notifier"
	"github.com/pkg/errors"
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
		return errors.New("store not loaded")
	}

	return r.remoteRepository.Start(ctx, onBootstrap)
}

func (r *ProxyStore[K]) Create(ctx context.Context, e K) (K, error) {
	if r.remoteRepository == nil {
		return e, errors.New("store not loaded")
	}

	if !e.GetId().IsEmpty() {
		return e, errors.New("entity already with id")
	}

	e, err := r.remoteRepository.Create(ctx, e)
	if err != nil {
		return e, errors.Wrap(err, "could not insert entity")
	}

	if r.Hydrate != nil {
		e, err = r.Hydrate(ctx, e)
		if err != nil {
			return e, errors.Wrap(err, "could not hydrate entity")
		}
	}
	r.localRepository.Create(ctx, e)
	r.notifier.Notify(ctx, Added, e)

	return e, nil
}

func (r *ProxyStore[K]) Read(ctx context.Context, id entity.Id) (e K, err error) {
	if r.remoteRepository == nil {
		return e, errors.New("store not loaded")
	}

	e, err = r.localRepository.Read(ctx, id)
	if err != nil && !errors.Is(err, entity.ErrEntityNotFound) {
		err = errors.Wrap(err, "could not read entity")
		return
	} else if err == nil {
		return
	}

	e, err = r.remoteRepository.Read(ctx, id)
	if err != nil && !errors.Is(err, entity.ErrEntityNotFound) {
		err = errors.Wrap(err, "could not read entity")
	} else if err == nil {
		r.localRepository.Create(ctx, e)
		r.notifier.Notify(ctx, Loaded, e)
	}

	return
}

func (r *ProxyStore[K]) Match(ctx context.Context, c criteria.Criteria) ([]K, error) {
	if r.remoteRepository == nil {
		return []K{}, errors.New("store not loaded")
	}

	return r.localRepository.Match(ctx, c)
}

func (r *ProxyStore[K]) MatchOne(ctx context.Context, c criteria.Criteria) (K, error) {
	if r.remoteRepository == nil {
		return *new(K), errors.New("store not loaded")
	}

	return r.localRepository.MatchOne(ctx, c)
}

func (r *ProxyStore[K]) ReadAll(ctx context.Context) ([]K, error) {
	if r.remoteRepository == nil {
		return []K{}, errors.New("store not loaded")
	}

	return r.localRepository.ReadAll(ctx)
}

func (r *ProxyStore[K]) Update(ctx context.Context, entity K) error {
	if r.remoteRepository == nil {
		return errors.New("store not loaded")
	}

	err := r.remoteRepository.Update(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "could not update entity")
	}

	r.localRepository.Update(ctx, entity)
	r.notifier.Notify(ctx, Updated, entity)

	return nil
}

func (r *ProxyStore[K]) Delete(ctx context.Context, entity K) error {
	if r.remoteRepository == nil {
		return errors.New("store not loaded")
	}

	err := r.remoteRepository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "could not delete entity")
	}

	r.localRepository.Delete(ctx, entity)
	r.notifier.Notify(ctx, Deleted, entity)

	return nil
}

func (r *ProxyStore[K]) Load(ctx context.Context, repo crudo.Repository[K]) error {
	r.remoteRepository = repo

	if r.localRepository != nil {
		return errors.New("Entities already loaded")
	}

	entities, err := r.remoteRepository.ReadAll(ctx)
	if err != nil {
		return errors.Wrap(err, "could not load Entities")
	}

	r.localRepository = inmemory.New(entities)
	for _, d := range entities {
		if r.Hydrate != nil {
			d, err = r.Hydrate(ctx, d)
			if err != nil {
				return errors.Wrap(err, "could not hydrate entity")
			}
		}

		r.notifier.Notify(ctx, Loaded, d)
	}

	return nil
}
func (r *ProxyStore[K]) Refresh(ctx context.Context) error {
	if r.remoteRepository == nil {
		return errors.New("store not loaded")
	}

	if r.RefreshPolicy == RefreshPolicyNone {
		return nil
	}

	remoteEntities, err := r.remoteRepository.ReadAll(ctx)
	if err != nil {
		return errors.Wrap(err, "could not load remote Entities")
	}

	localEntities, err := r.localRepository.ReadAll(ctx)
	if err != nil {
		return errors.Wrap(err, "could not load local Entities")
	}

	if r.RefreshPolicy == RefreshPolicyReadAll || r.RefreshPolicy == RefreshPolicyReadWriteAll {
		for _, d := range remoteEntities {
			if !entity.Contains(localEntities, d) {
				r.localRepository.Create(ctx, d)
				r.notifier.Notify(ctx, Loaded, d)
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
				r.localRepository.Delete(ctx, d)
				r.notifier.Notify(ctx, Unloaded, d)
			}
		}
	}

	return nil
}
