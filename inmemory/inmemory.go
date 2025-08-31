package inmemory

import (
	"context"
	"fmt"

	"sync"

	"github.com/davfer/archit/patterns/opts"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
	"github.com/google/uuid"
)

type Repository[K entity.Entity] struct {
	Collection []K
	lock       *sync.Mutex
	policy     Policy[K]
	idStrategy IdStrategy[K]
}

func NewRepository[K entity.Entity](c []K, o ...opts.Opt[Repository[K]]) *Repository[K] {
	r := opts.New[Repository[K]](o...)

	r.Collection = c
	r.lock = &sync.Mutex{}

	return &r
}

func (r *Repository[K]) Start(ctx context.Context, onBootstrap func(ctx context.Context) error) error {
	return onBootstrap(ctx)
}

func (r *Repository[K]) Create(ctx context.Context, e K) (K, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if entity.Contains(r.Collection, e) {
		return e, entity.ErrEntityAlreadyExists
	}

	if ee, ok := entity.Entity(e).(entity.EventfulEntity); ok {
		err := ee.PreCreate()
		if err != nil {
			return e, fmt.Errorf("error pre creating entity: %w", err)
		}
	}

	if r.idStrategy != nil {
		id := r.idStrategy.Generate(e)

		err := e.SetId(id)
		if err != nil {
			return e, fmt.Errorf("error setting generated entity id: %w", err)
		}
	} else if e.GetId().IsEmpty() {
		err := e.SetId(entity.NewIdFromString(uuid.New().String()))
		if err != nil {
			return e, fmt.Errorf("error setting entity id: %w", err)
		}
	}

	if r.policy != nil {
		c, err := r.policy.ApplyCreate(ctx, e, r.Collection)
		if err != nil {
			return e, err
		}

		r.Collection = c
	} else { // Nil policy, just append
		r.Collection = append(r.Collection, e)
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

func (r *Repository[K]) Match(ctx context.Context, c specification.Criteria) ([]K, error) {
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

func (r *Repository[K]) MatchOne(ctx context.Context, c specification.Criteria) (k K, err error) {
	ks, _ := r.Match(ctx, c)
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
