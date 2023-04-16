package event

import "context"

type TopicEvent struct {
	Event     Event
	TimeStamp int64
	Topic     string
}

type Bus interface {
	Listener() Listener
	Publish(ctx context.Context, topic string, event Event) error
}

type Listener interface {
	Subscribe(ctx context.Context, topics ...string) error
	Unsubscribe(ctx context.Context, topics ...string) error
	UnsubscribeAll(ctx context.Context) error
	Listen(ctx context.Context) (<-chan TopicEvent, error)
}
