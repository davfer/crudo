package inmemory

import (
	"github.com/davfer/crudo/entity"
	"sync"
)

type Opt[K entity.Entity] func(repository *Repository[K])

func New[K entity.Entity](data []K, opts ...Opt[K]) *Repository[K] {
	r := &Repository[K]{Collection: data, lock: &sync.Mutex{}}
	for _, o := range opts {
		o(r)
	}
	return r
}

func WithPolicy[K entity.Entity](p Policy[K]) Opt[K] {
	return func(s *Repository[K]) {
		s.policy = p
	}
}

func WithIdStrategy[K entity.Entity](i IdStrategy[K]) Opt[K] {
	return func(s *Repository[K]) {
		s.idStrategy = i
	}
}
