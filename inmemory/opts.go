package inmemory

import (
	"github.com/davfer/archit/patterns/opts"
	"github.com/davfer/crudo/entity"
)

func WithPolicy[K entity.Entity](p Policy[K]) opts.Opt[Repository[K]] {
	return func(s Repository[K]) Repository[K] {
		s.policy = p
		return s
	}
}

func WithIdStrategy[K entity.Entity](i IdStrategy[K]) opts.Opt[Repository[K]] {
	return func(s Repository[K]) Repository[K] {
		s.idStrategy = i
		return s
	}
}
