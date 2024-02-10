package notifier

import "context"

type TopicCallbackNotifier[K any] struct {
	notifier *TopicNotifier[K]
}

func NewTopicCallbackNotifier[K any](topics []string) *TopicCallbackNotifier[K] {
	t := &TopicCallbackNotifier[K]{
		notifier: NewTopicNotifier[K](topics),
	}

	return t
}

func (t *TopicCallbackNotifier[K]) Attach(topic string, observer ObserverCallback[K]) error {
	cb := &topicSubscriber[K]{
		cb: observer,
	}

	return t.notifier.Attach(topic, cb)
}

func (t *TopicCallbackNotifier[K]) Notify(ctx context.Context, topic string, event K) error {
	return t.notifier.Notify(ctx, topic, event)
}

type topicSubscriber[K any] struct {
	cb ObserverCallback[K]
}

func (t *topicSubscriber[K]) Handle(ctx context.Context, event K) error {
	return t.cb(ctx, event)
}
