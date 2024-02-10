package inmemory

import (
	"context"
	"github.com/davfer/crudo/criteria"
	"github.com/davfer/crudo/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync"
)

type Repository[K entity.Entity] struct {
	Collection []K
	lock       *sync.Mutex
	policy     Policy[K]
	idStrategy IdStrategy[K]
}

func (r *Repository[K]) Start(ctx context.Context, onBootstrap func(ctx context.Context) error) error {
	return nil
}

func (r *Repository[K]) Create(ctx context.Context, e K) (K, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, i := range r.Collection {
		if i.GetId() == e.GetId() {
			return e, entity.ErrEntityAlreadyExists
		}
	}

	err := e.PreCreate()
	if err != nil {
		return e, errors.Wrap(err, "error pre creating entity")
	}

	if r.policy != nil {
		r.Collection, err = r.policy.ApplyCreate(ctx, e, r.Collection)
		if err != nil {
			return e, err
		}
	} else { // Nil policy, just append
		r.Collection = append(r.Collection, e)
	}

	if e.GetId().IsEmpty() {
		err = e.SetId(entity.NewIdFromString(uuid.New().String()))
		if err != nil {
			return e, errors.Wrap(err, "error setting entity id")
		}
	}

	return e, nil
}

func (r *Repository[K]) Read(ctx context.Context, id entity.Id) (e K, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, i := range r.Collection {
		if i.GetId() == id {
			e = i
			return
		}
	}

	err = entity.ErrEntityNotFound
	return
}

func (r *Repository[K]) Match(ctx context.Context, c criteria.Criteria) ([]K, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var result []K
	for _, e := range r.Collection {
		if c.IsSatisfiedBy(e) {
			result = append(result, e)
		}
	}

	return result, nil
}

func (r *Repository[K]) MatchOne(ctx context.Context, c criteria.Criteria) (k K, err error) {
	ks, err := r.Match(ctx, c)
	if err != nil {
		return k, err
	}
	if len(ks) == 0 {
		return k, entity.ErrEntityNotFound
	}

	k = ks[0]
	return
}

func (r *Repository[K]) ReadAll(ctx context.Context) ([]K, error) {
	return r.Collection, nil
}

func (r *Repository[K]) Update(ctx context.Context, entity K) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, e := range r.Collection {
		if e.GetId() == entity.GetId() {
			r.Collection[i] = entity
			break
		}
	}

	return nil
}

func (r *Repository[K]) Delete(ctx context.Context, entity K) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, e := range r.Collection {
		if e.GetId() == entity.GetId() {
			r.Collection = append(r.Collection[:i], r.Collection[i+1:]...)
			break
		}
	}

	return nil
}
