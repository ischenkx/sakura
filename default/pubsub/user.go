package pubsub

import (
	"github.com/RomanIschenko/notify/default/pubsub/subscription"
	"sync"
)

type user struct {
	clients       map[string]*client
	subscriptions subscription.List
	data          *sync.Map
}

func (u *user) Add(c *client) {
	u.clients[c.id] = c
}

func (u *user) Del(id string) {
	delete(u.clients, id)
}

func (u *user) Subscribe(id string, ts int64, f func(*client)) error {
	err := u.subscriptions.Add(id, ts)
	if err == nil {
		for _, c := range u.clients {
			if err := c.subscriptions.Add(id, ts); err == nil && f != nil {
				f(c)
			}
		}
		return nil
	}
	return err
}

func (u *user) Unsubscribe(id string, ts int64, f func(*client)) error {
	err := u.subscriptions.Delete(id, ts)
	if err == nil {
		for _, c := range u.clients {
			if err := c.subscriptions.Delete(id, ts); err == nil && f != nil {
				f(c)
			}
		}
		return nil
	}
	return err
}

func newUser() *user {
	return &user{
		data:          &sync.Map{},
		clients:       map[string]*client{},
		subscriptions: subscription.NewList(),
	}
}
