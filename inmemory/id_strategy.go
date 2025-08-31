package inmemory

import (
	"github.com/davfer/crudo/entity"
	"github.com/google/uuid"
)

type IdStrategy[K entity.Entity] interface {
	Generate(k K) entity.ID
}

type UuidIdStrategy[K entity.Entity] struct{}

func (d UuidIdStrategy[K]) Generate(k K) entity.ID {
	return entity.NewIDFromString(uuid.New().String())
}
