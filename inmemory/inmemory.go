package inmemory

import (
	"context"
	"github.com/davfer/crudo/criteria"
	"github.com/davfer/crudo/entity"
	"github.com/pkg/errors"
)

type Repository[K entity.Entity] struct {
	Collection []K
}

func NewInMemoryRepository[K entity.Entity](data []K) *Repository[K] {
	return &Repository[K]{Collection: data}
}

func (r *Repository[K]) Start(ctx context.Context, onBootstrap func(ctx context.Context) error) error {
	return nil
}

func (r *Repository[K]) Create(ctx context.Context, e K) (entity.Id, error) {
	for _, i := range r.Collection {
		if i.GetId() == e.GetId() {
			return "", entity.ErrEntityAlreadyExists
		}
	}

	err := e.PreCreate()
	if err != nil {
		return "", errors.Wrap(err, "error pre creating entity")
	}

	r.Collection = append(r.Collection, e)

	return e.GetId(), nil
}

func (r *Repository[K]) Read(ctx context.Context, id entity.Id) (K, error) {
	var e K

	for _, i := range r.Collection {
		if i.GetId() == id {
			return e, nil
		}
	}

	return e, entity.ErrEntityNotFound
}

func (r *Repository[K]) Match(ctx context.Context, c criteria.Criteria) ([]K, error) {
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
	return nil
}

func (r *Repository[K]) Delete(ctx context.Context, entity K) error {
	for i, e := range r.Collection {
		if e.GetId() == entity.GetId() {
			r.Collection = append(r.Collection[:i], r.Collection[i+1:]...)
			break
		}
	}

	return nil
}
