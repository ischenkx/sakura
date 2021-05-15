package user

import (
	"github.com/RomanIschenko/notify/internal/pubsub/internal/client"
	"github.com/RomanIschenko/notify/internal/pubsub/subscription"
)

type User struct {
	clients       map[string]*client.Client
	subscriptions subscription.List
}

func (u *User) Add(c *client.Client) {
	u.clients[c.ID()] = c
}

func (u *User) Del(id string) {
	delete(u.clients, id)
}

func (u *User) Subscribe(id string, ts int64, f func(c *client.Client)) error {
	err := u.subscriptions.Add(id, ts)
	if err == nil {
		for _, c := range u.clients {
			if err := c.Subscriptions().Add(id, ts); err == nil && f != nil {
				f(c)
			}
		}
		return nil
	}
	return err
}

func (u *User) Unsubscribe(id string, ts int64, f func(client *client.Client)) error {
	err := u.subscriptions.Delete(id, ts)
	if err == nil {
		for _, c := range u.clients {
			if err := c.Subscriptions().Delete(id, ts); err == nil && f != nil {
				f(c)
			}
		}
		return nil
	}
	return err
}

func (u *User) Subscriptions() subscription.List {
	return u.subscriptions
}

func (u *User) Clients() map[string]*client.Client {
	return u.clients
}

func New() *User {
	return &User{
		clients:       map[string]*client.Client{},
		subscriptions: subscription.NewList(),
	}
}
