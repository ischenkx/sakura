package user

import (
	"github.com/RomanIschenko/notify/pubsub/internal/client"
	"github.com/RomanIschenko/notify/pubsub/internal/subscription"
)

type User struct {
	descriptor uint64
	clients map[string]*client.Client
	subscriptions subscription.List
}

func (u *User) Subscribe(t string, seq int64, f func(*client.Client)) {
	if err := u.subscriptions.Add(t, seq); err == nil {
		for _, c := range u.clients {
			if err := c.Subscribe(t, seq); err == nil {
				f(c)
			}
		}
	}
}

func (u *User) Unsubscribe(t string, seq int64, f func(*client.Client)) {
	if err := u.subscriptions.Remove(t, seq); err == nil {
		for _, c := range u.clients {
			if err := c.Unsubscribe(t, seq); err == nil {
				f(c)
			}
		}
	}
}

func (u *User) Subscriptions() map[string]subscription.Subscription {
	return u.subscriptions.Map()
}

func (u *User) Clients() map[string]*client.Client {
	return u.clients
}

func (u *User) Descriptor() uint64 {
	return u.descriptor
}

func New(desc uint64) *User {
	return &User{
		clients: map[string]*client.Client{},
		descriptor:    desc,
		subscriptions: subscription.NewList(),
	}
}