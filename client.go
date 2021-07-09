package swirl

import (
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"time"
)

type localClient struct {
	id  string
	app *App
}

func (c localClient) Emit(name string, options ...interface{}) {
	var args []interface{}
	var metaInfo interface{}
	var timeStamp int64

	for _, opt := range options {
		switch o := opt.(type) {
		case MetaInfo:
			metaInfo = o
		case TimeStamp:
			timeStamp = time.Time(o).UnixNano()
		case Args:
			args = o
		}
	}
	data, err := c.app.emitter.EncodeRawData(name, args)
	if err != nil {
		c.app.events.callError(EncodingError{
			Reason: err,
			EventOptions: EventOptions{
				Name:      name,
				Args:      args,
				TimeStamp: timeStamp,
				MetaInfo:  metaInfo,
			},
		})
		return
	}
	c.app.pubsub.SendToClient(c.id, message.New(data))
	c.app.events.callEmit(EmitOptions{
		Clients: []string{c.id},
		EventOptions: EventOptions{
			Name:      name,
			Args:      args,
			TimeStamp: timeStamp,
			MetaInfo:  metaInfo,
		},
	})
}

func (c localClient) Subscribe(options SubscribeOptions) {
	cl := ChangeLog{TimeStamp: options.TimeStamp}
	for _, topic := range options.Topics {
		if locCl, err := c.app.pubsub.Subscribe(c.id, topic, options.TimeStamp); err == nil {
			cl.Merge(locCl)
			c.app.events.callClientSubscribe(c, topic, options.TimeStamp)
		}
	}
	c.app.events.callChange(cl)
}

func (c localClient) Unsubscribe(options UnsubscribeOptions) {
	cl := ChangeLog{TimeStamp: options.TimeStamp}
	for _, topic := range options.Topics {
		if locCl, err := c.app.pubsub.Unsubscribe(c.id, topic, options.TimeStamp); err == nil {
			cl.Merge(locCl)
			c.app.events.callClientUnsubscribe(c, topic, options.TimeStamp)
		}
	}
	c.app.events.callChange(cl)
}

func (c localClient) Disconnect(options DisconnectOptions) {
	defer c.app.events.callDisconnect(c)
	cl := ChangeLog{TimeStamp: options.TimeStamp}
	cl.Merge(c.app.pubsub.Disconnect(c.id, options.TimeStamp))
	cl.Merge(c.app.pubsub.Users().Delete(c.userID(), c.id, options.TimeStamp))
	c.app.events.callChange(cl)

}

func (c localClient) Active() bool {
	return c.app.pubsub.IsActive(c.id)
}

func (c localClient) userID() string {
	return c.app.pubsub.Users().UserByClient(c.id)
}

func (c localClient) User() User {
	return c.app.User(c.userID())
}

func (c localClient) ID() string {
	return c.id
}

func (c localClient) Subscriptions() SubscriptionList {
	return variadicSubList{
		count: func(list SubscriptionList) int {
			return c.app.pubsub.CountSubscriptions(c.id)
		},
		array: func(list SubscriptionList) []Subscription {
			return c.app.pubsub.Subscriptions(c.id)
		},
	}
}

func (c localClient) Events() ClientEvents {
	return c.app.events.forClient(c.id)
}

type Client interface {
	ID() string
	Emit(string, ...interface{})
	Subscribe(options SubscribeOptions)
	Unsubscribe(options UnsubscribeOptions)
	Disconnect(options DisconnectOptions)
	User() User
	Active() bool
	Events() ClientEvents
	Subscriptions() SubscriptionList
}
