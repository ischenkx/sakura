package pubsub

import (
	"github.com/RomanIschenko/notify/default/pubsub/subscription"
	"sync"
)

type client struct {
	id, user      string
	meta          *sync.Map
	subscriptions subscription.List
}

func (c *client) ID() string {
	return c.id
}

func (c *client) User() string {
	return c.user
}

func (c *client) Data() *sync.Map {
	return c.meta
}

func newClient(id, user string) *client {
	return &client{
		id:            id,
		user:          user,
		meta:          &sync.Map{},
		subscriptions: subscription.NewList(),
	}
}
