package topic

import (
	"errors"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/namespace"
)

type Topic struct {
	cfg   namespace.Config
	subs  map[*pubsub.Client]struct{}
	users map[string]int
}

func (t Topic) Add(c *pubsub.Client) error {
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

func (t Topic) Del(c *pubsub.Client) (int, error) {
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

func (t Topic) Subscribers() map[*pubsub.Client]struct{} {
	return t.subs
}

func New(cfg namespace.Config) Topic {
	return Topic{
		cfg:   cfg,
		subs: map[*pubsub.Client]struct{}{},
		users: map[string]int{},
	}
}
