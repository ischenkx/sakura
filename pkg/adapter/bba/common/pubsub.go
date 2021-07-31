package common

import "context"

type PubSub interface {
	Subscribe(context.Context, ...string) error
	Unsubscribe(context.Context, ...string) error
	UnsubscribeAll(context.Context) error
	Publish(context.Context, string, Event) error
	Channel() <-chan Event
	Close()
}
