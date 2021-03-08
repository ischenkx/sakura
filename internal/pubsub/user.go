package pubsub

import "github.com/RomanIschenko/notify/internal/pubsub/internal/subscription"

type user struct {
	clients map[string]*client
	subscriptions subscription.List
}

func (u *user) Add(c *client) {
	u.clients[c.id] = c
}

func (u *user) Del(id string) {
	delete(u.clients, id)
}

func (u *user) Subscribe(id string, ts int64, f func(*client)) {
	if err := u.subscriptions.Add(id, ts); err == nil {
		for _, c := range u.clients {
			if err := c.subscriptions.Add(id, ts); err == nil && f != nil {
				f(c)
			}
		}
	}
}

func (u *user) Unsubscribe(id string, ts int64, f func(*client)) {
	if err := u.subscriptions.Delete(id, ts); err == nil {
		for _, c := range u.clients {
			if err := c.subscriptions.Delete(id, ts); err == nil && f != nil {
				f(c)
			}
		}
	}
}

func newUser() *user {
	return &user{
		clients: map[string]*client{},
		subscriptions: subscription.NewList(),
	}
}