package broker

import "context"

type Message[T any] struct {
	Channel string
	Data    T
}

type PubSub[T any] interface {
	Subscribe(ctx context.Context, channels ...string) error
	Unsubscribe(ctx context.Context, channels ...string) error
	Channel(ctx context.Context) (<-chan Message[T], error)
	Clear(ctx context.Context) error
}

type Broker[T any] interface {
	Push(ctx context.Context, channel string, message T) error
	PubSub() PubSub[T]
}
