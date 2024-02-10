package notifier

import (
	"context"
	"github.com/pkg/errors"
)

type Observer[K any] interface {
	Handle(context.Context, K) error
}

type ObserverCallback[K any] func(context.Context, K) error

type Notifier[K any] struct {
	observers []Observer[K]
}

func NewNotifier[K any]() *Notifier[K] {
	n := &Notifier[K]{
		observers: []Observer[K]{},
	}

	return n
}

func (n *Notifier[K]) Attach(observer Observer[K]) error {
	n.observers = append(n.observers, observer)

	return nil
}

func (n *Notifier[K]) Detach(observer Observer[K]) error {
	for i, h := range n.observers {
		if h == observer {
			n.observers = append(n.observers[:i], n.observers[i+1:]...)
			break
		}
	}

	return nil
}

func (n *Notifier[K]) Notify(ctx context.Context, event K) error {
	var err error
	for _, d := range n.observers {
		err = d.Handle(ctx, event)
		if err != nil {
			err = errors.Wrap(err, "failed to notify observer")
		}
	}

	return err
}
