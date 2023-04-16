package sakura

import (
	"sakura/event"
	"sakura/subscription"
)

type Sakura struct {
	subscriptions subscription.Storage
	events        event.Bus
}

func New(subscriptions subscription.Storage, events event.Bus) *Sakura {
	return &Sakura{
		subscriptions: subscriptions,
		events:        events,
	}
}

func (sakura *Sakura) User(id string) User {
	return User{
		id:     id,
		sakura: sakura,
	}
}

func (sakura *Sakura) Topic(id string) Topic {
	return Topic{
		id:     id,
		sakura: sakura,
	}
}

func (sakura *Sakura) Events() event.Bus {
	return sakura.events
}
