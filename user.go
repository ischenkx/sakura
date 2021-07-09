package swirl

import (
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"time"
)

type localUser struct {
	id  string
	app *App
}

func (u localUser) Emit(name string, options ...interface{}) {

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

	if timeStamp == 0 {
		timeStamp = time.Now().UnixNano()
	}

	data, err := u.app.emitter.EncodeRawData(name, args)
	if err != nil {
		u.app.events.callError(EncodingError{
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

	u.app.pubsub.Users().Iter(u.id, func(client string) {
		u.app.pubsub.SendToClient(client, message.New(data))
	})

	u.app.events.callEmit(EmitOptions{
		Users: []string{u.id},
		EventOptions: EventOptions{
			Name:      name,
			Args:      args,
			TimeStamp: timeStamp,
			MetaInfo:  metaInfo,
		},
	})
}

func (u localUser) Subscribe(options SubscribeOptions) {
	var subscribedTopics []string
	for _, topic := range options.Topics {
		if err := u.app.pubsub.Users().Subscribe(u.id, topic, options.TimeStamp); err == nil {
			subscribedTopics = append(subscribedTopics, topic)
		}
	}
	u.app.pubsub.Users().Iter(u.id, func(client string) {
		for _, topic := range subscribedTopics {
			u.app.pubsub.Subscribe(client, topic, options.TimeStamp)
		}
	})

	for _, topic := range subscribedTopics {
		u.app.events.callUserSubscribe(u, topic, options.TimeStamp)
	}
}

func (u localUser) Active() (isActive bool) {
	return u.app.pubsub.Users().IsActive(u.id)
}

func (u localUser) Unsubscribe(options UnsubscribeOptions) {
	var unsubscribedTopics []string

	for _, topic := range options.Topics {
		if err := u.app.pubsub.Users().Unsubscribe(u.id, topic, options.TimeStamp); err == nil {
			unsubscribedTopics = append(unsubscribedTopics, topic)
		}
	}

	u.app.pubsub.Users().Iter(u.id, func(client string) {
		for _, topic := range unsubscribedTopics {
			u.app.pubsub.Unsubscribe(client, topic, options.TimeStamp)
		}
	})

	for _, topic := range unsubscribedTopics {
		u.app.events.callUserUnsubscribe(u, topic, options.TimeStamp)
	}
}

func (u localUser) Disconnect(options DisconnectOptions) {
	u.app.pubsub.Users().Iter(u.id, func(client string) {
		u.app.pubsub.Disconnect(client, options.TimeStamp)
	})
}

func (u localUser) Clients() []string {
	clients := make([]string, 0)
	u.app.pubsub.Users().Iter(u.id, func(id string) {
		clients = append(clients, id)
	})
	return clients
}

func (u localUser) Subscriptions() SubscriptionList {
	return variadicSubList{
		count: func(list SubscriptionList) int {
			return u.app.pubsub.Users().CountSubscriptions(u.id)
		},
		array: func(list SubscriptionList) []Subscription {
			subs := make([]Subscription, 0, list.Count())
			u.app.pubsub.Users().IterSubscriptions(u.id, func(s Subscription) {
				subs = append(subs, s)
			})
			return subs
		},
	}
}

func (u localUser) Events() UserEvents {
	return u.app.events.forUser(u.id)
}

func (u localUser) ID() string {
	return u.id
}

type User interface {
	ID() string
	Active() bool
	Emit(name string, data ...interface{})
	Subscribe(options SubscribeOptions)
	Unsubscribe(options UnsubscribeOptions)
	Disconnect(options DisconnectOptions)
	Clients() []string
	Events() UserEvents
	Subscriptions() SubscriptionList
}
