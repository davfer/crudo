package store

import (
	"context"
	"github.com/davfer/crudo"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/notifier"
)

type Store[K entity.Entity] interface {
	crudo.Repository[K]
	// Load hydrates the store with the given repository
	Load(context.Context, crudo.Repository[K]) error
	// Refresh updates the store with the latest data from the repository
	Refresh(context.Context) error
	// OnHydrate sets the hydrate function to be called when an entity is loaded
	OnHydrate(hydrateFunc HydrateFunc[K])
	// On attaches an observer to the store
	On(string, notifier.ObserverCallback[K]) error
}
