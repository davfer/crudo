package notifier

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"
)

type topicObserver[K any] struct {
	topic    string
	observer Observer[K]
}

type TopicNotifier[K any] struct {
	observers []topicObserver[K]
	topics    []string
}

func NewTopicNotifier[K any](topics []string) *TopicNotifier[K] {
	return &TopicNotifier[K]{
		observers: []topicObserver[K]{},
		topics:    topics,
	}
}

func (n *TopicNotifier[K]) Attach(topic string, observer Observer[K]) error {
	if !slices.Contains(n.topics, topic) {
		return fmt.Errorf("topic %s not found", topic)
	}

	n.observers = append(n.observers, topicObserver[K]{
		topic:    topic,
		observer: observer,
	})

	return nil
}

func (n *TopicNotifier[K]) Detach(observer Observer[K]) error {
	for i, h := range n.observers {
		if h.observer == observer {
			n.observers = append(n.observers[:i], n.observers[i+1:]...)
			break
		}
	}

	return nil
}

func (n *TopicNotifier[K]) Notify(ctx context.Context, topic string, event K) error {
	var err error
	for _, o := range n.observers {
		if o.topic == topic {
			err = o.observer.Handle(ctx, event)
			if err != nil {
				err = fmt.Errorf("failed to handle event: %w", err)
			}
		}
	}

	return err
}
