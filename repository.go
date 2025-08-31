package crudo

import (
	"context"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
)

type Repository[K entity.Entity] interface {
	Start(ctx context.Context, onBootstrap func(context.Context) error) error
	QueryRepository[K]
	WriteRepository[K]
}

type QueryRepository[K entity.Entity] interface {
	Read(context.Context, entity.ID) (K, error)
	ReadAll(context.Context) ([]K, error)
	Match(context.Context, specification.Criteria) ([]K, error)
	MatchOne(context.Context, specification.Criteria) (K, error)
}

type WriteRepository[K entity.Entity] interface {
	Create(context.Context, K) (K, error)
	Update(context.Context, K) error
	Delete(context.Context, K) error
}
