package broker

import "context"

type PubSub[Message any] interface {
	Subscribe(ctx context.Context, channels ...string) error
	Unsubscribe(ctx context.Context, channels ...string) error
	Channel(ctx context.Context) (<-chan Message, error)
	Clear(ctx context.Context) error
}

type Broker[Message any] interface {
	Push(ctx context.Context, channel string, message Message) error
	PubSub(ctx context.Context) PubSub[Message]
}
