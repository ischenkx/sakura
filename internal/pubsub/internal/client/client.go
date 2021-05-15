package client

import (
	subscription2 "github.com/RomanIschenko/notify/internal/pubsub/subscription"
	"sync/atomic"
)

type Client struct {
	id, user      string
	subscriptions subscription2.List

	state int64
}

func (c *Client) Invalidate() {
	atomic.StoreInt64(&c.state, 1)
}

func (c *Client) Valid() bool {
	return atomic.LoadInt64(&c.state) == 0
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) User() string {
	return c.user
}

func (c *Client) Subscriptions() subscription2.List {
	return c.subscriptions
}

func New(id, user string) *Client {
	return &Client{
		id:            id,
		user:          user,
		subscriptions: subscription2.NewList(),
	}
}
