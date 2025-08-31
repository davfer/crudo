package inmemory

import (
	"context"

	"github.com/davfer/crudo/entity"
)

type Policy[K entity.Entity] interface {
	ApplyCreate(ctx context.Context, e K, col []K) ([]K, error)
}

type PolicyMru[K entity.Entity] struct {
	Capacity int
}

func (p PolicyMru[K]) ApplyCreate(ctx context.Context, e K, col []K) ([]K, error) {
	if len(col) >= p.Capacity {
		col = col[1:]
	}

	return append(col, e), nil
}

type PolicyLru[K entity.Entity] struct {
	Capacity int
}

func (p PolicyLru[K]) ApplyCreate(ctx context.Context, e K, col []K) ([]K, error) {
	if len(col) >= p.Capacity {
		return col, nil
	}

	return append(col, e), nil
}
