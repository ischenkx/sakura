package pubsub

import (
	"github.com/RomanIschenko/notify/pubsub/internal/client"
	"sync"
)

type Client struct {
	raw *client.Client
}

func (c Client) User() string {
	return c.raw.User()
}

func (c Client) ID() string {
	return c.raw.ID()
}

func (c Client) Meta() *sync.Map {
	return c.raw.Meta()
}
