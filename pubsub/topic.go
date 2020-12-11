package pubsub

import (
	"errors"
	"github.com/RomanIschenko/notify/pubsub/namespace"
)

type topic struct {
	cfg   namespace.Config
	subs  map[*Client]struct{}
	users map[string]int
}

func (t topic) add(c *Client) error {
	if t.cfg.MaxUsers != namespace.Any {
		userID := c.UserID()
		if _, ok := t.users[userID]; !ok {
			if len(t.users)+1 > t.cfg.MaxUsers {
				return errors.New("cfg.MaxUsers is overflowed")
			}
		}
		counter := t.users[userID]
		t.users[userID] = counter + 1
	}
	if t.cfg.MaxClients != namespace.Any {
		if _, ok := t.subs[c]; !ok {
			if len(t.subs) + 1 > t.cfg.MaxClients {
				return errors.New("cfg.MaxClients is overflowed")
			}
		}
	}
	t.subs[c] = struct{}{}
	return nil
}

func (t topic) del(c *Client) (int, error) {
	if _, ok := t.subs[c]; ok {
		delete(t.subs, c)
		if t.cfg.MaxUsers != namespace.Any {
			userID := c.UserID()
			if counter, ok := t.users[userID]; ok {
				counter -= 1
				if counter <= 0 {
					delete(t.users, userID)
				} else {
					t.users[userID] = counter
				}
			}
		}
		return len(t.subs), nil
	}
	return -1, errors.New("no such client")
}

func (t topic) subscribers() map[*Client]struct{} {
	return t.subs
}

func newTopic(cfg namespace.Config) topic {
	return topic{
		cfg:   cfg,
		subs: map[*Client]struct{}{},
		users: map[string]int{},
	}
}
