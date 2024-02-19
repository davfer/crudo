package crudo

import (
	"context"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
)

type Repository[K entity.Entity] interface {
	Start(ctx context.Context, onBootstrap func(context.Context) error) error

	Create(context.Context, K) (K, error)
	Read(context.Context, entity.Id) (K, error)
	ReadAll(context.Context) ([]K, error)
	Update(context.Context, K) error
	Delete(context.Context, K) error

	Match(context.Context, specification.Criteria) ([]K, error)
	MatchOne(context.Context, specification.Criteria) (K, error)
}
