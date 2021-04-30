package pubsub

import (
	"github.com/RomanIschenko/notify/default/pubsub/subscription"
	"github.com/RomanIschenko/notify/pubsub"
	"sync"
)

type client struct {
	id, user      string
	data          *sync.Map
	userData      *sync.Map
	subscriptions subscription.List
}

func (c *client) ID() string {
	return c.id
}

func (c *client) User() string {
	return c.user
}

func (c *client) Data() pubsub.LocalStorage {
	return c.data
}

func (c *client) UserData() pubsub.LocalStorage {
	return c.data
}

func newClient(id, user string) *client {
	return &client{
		id:            id,
		user:          user,
		data:          &sync.Map{},
		subscriptions: subscription.NewList(),
	}
}
