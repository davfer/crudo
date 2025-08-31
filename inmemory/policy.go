package inmemory

import (
	"context"

	"github.com/davfer/crudo/entity"
)

type Policy[K entity.Entity] interface {
	ApplyCreate(ctx context.Context, e K, col []K) ([]K, error)
}

type PolicyMRU[K entity.Entity] struct {
	Capacity int
}

func (p PolicyMRU[K]) ApplyCreate(ctx context.Context, e K, col []K) ([]K, error) {
	if len(col) >= p.Capacity {
		col = col[1:]
	}

	return append(col, e), nil
}

type PolicyLRU[K entity.Entity] struct {
	Capacity int
}

func (p PolicyLRU[K]) ApplyCreate(ctx context.Context, e K, col []K) ([]K, error) {
	if len(col) >= p.Capacity {
		return col, nil
	}

	return append(col, e), nil
}
