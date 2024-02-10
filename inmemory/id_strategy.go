package inmemory

import (
	"github.com/davfer/crudo/entity"
	"github.com/google/uuid"
)

type IdStrategy[K entity.Entity] interface {
	Generate(k K) entity.Id
}

type UuidIdStrategy[K entity.Entity] struct{}

func (d UuidIdStrategy[K]) Generate(k K) entity.Id {
	return entity.NewIdFromString(uuid.New().String())
}
