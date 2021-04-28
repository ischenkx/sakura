package pubsub

import (
	"context"
)

type Engine interface {
	Clean() ([]Client, ChangeLog)
	Connect(opts ConnectOptions) (Client, ChangeLog, bool, error)
	Inactivate(id string, ts int64) (Client, ChangeLog, error)
	Disconnect(opts DisconnectOptions) ([]Client, ChangeLog)
	Subscribe(opts SubscribeOptions) ChangeLog
	Unsubscribe(opts UnsubscribeOptions) ChangeLog
	Publish(opts PublishOptions)
	Start(ctx context.Context)
	IsSubscribed(client string, topic string) (bool, error)
	IsUserSubscribed(user string, topic string) (bool, error)
	TopicSubscribers(topic string) ([]string, error)
	ClientSubscriptions(id string) ([]string, error)
	UserSubscriptions(id string) ([]string, error)
}
