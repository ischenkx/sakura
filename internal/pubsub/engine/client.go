package engine

import (
	"github.com/ischenkx/swirl/internal/pubsub/session"
	"github.com/ischenkx/swirl/internal/pubsub/subscription"
)

type ClientInfo struct {
	Subscriptions *subscription.List
	ID            string
}

type client struct {
	subscriptions *subscription.List
	session       *session.Session
}

func newClient(id string, ts int64) *client {
	return &client{
		subscriptions: subscription.NewList(),
		session:       session.New(id, nil, ts),
	}
}
