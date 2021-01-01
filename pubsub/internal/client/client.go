package client

import (
	"github.com/RomanIschenko/notify/pubsub/internal/subscription"
	"sync"
)

type Client struct {
	descriptor uint64
	id, user string
	meta *sync.Map
	subscriptions subscription.List
}

func (c *Client) Subscribe(t string, seq int64) error {
	return c.subscriptions.Add(t, seq)
}

func (c *Client) Unsubscribe(t string, seq int64) error {
	return c.subscriptions.Remove(t, seq)
}

func (c *Client) Subscriptions() map[string]subscription.Subscription {
	return c.subscriptions.Map()
}

func (c *Client) Descriptor() uint64 {
	return c.descriptor
}

func (c *Client) User() string {
	return c.user
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Meta() *sync.Map {
	return c.meta
}

func New(desc uint64, id, user string) *Client {
	return &Client{
		descriptor:    desc,
		id: 		   id,
		user:          user,
		meta: 		   &sync.Map{},
		subscriptions: subscription.NewList(),
	}
}


