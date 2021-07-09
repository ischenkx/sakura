package engine

import (
	"github.com/ischenkx/swirl/internal/pubsub/session"
	"github.com/ischenkx/swirl/internal/pubsub/subscription"
)

type memberInfo struct {
	Ok        bool
	TimeStamp int64
}

type topic struct {
	subscriptions *subscription.List
	clients       map[string]memberInfo
	group *session.Group
}

func (u *topic) addClient(id string, ts int64) {
	info := u.clients[id]
	if info.TimeStamp > ts {
		return
	}
	info.Ok = true
	info.TimeStamp = ts
	u.clients[id] = info
}

func (u *topic) deleteClient(id string, forced bool, ts int64) {
	info := u.clients[id]

	if forced {
		delete(u.clients, id)
		return
	}

	if info.TimeStamp > ts {
		return
	}

	info.Ok = false
	info.TimeStamp = ts
	u.clients[id] = info
}

func newTopic() *topic {
	return &topic{
		subscriptions: subscription.NewList(),
		clients:       map[string]memberInfo{},
		group: session.NewGroup(),
	}
}
