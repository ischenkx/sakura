package pubsub

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify/internal/pubsub/broadcaster"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/client"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/topic"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/user"
	"github.com/RomanIschenko/notify/internal/pubsub/protocol"
	"github.com/RomanIschenko/notify/pkg/default/batchproto"
	"sync"
	"time"
)

type Config struct {
	InvalidationTime time.Duration
	CleanInterval    time.Duration
	ProtoProvider    protocol.Provider
}

type PubSub struct {
	clients         map[string]*client.Client
	topics          map[string]*topic.Topic
	users           map[string]*user.User
	inactiveClients map[string]int64

	broadcaster *broadcaster.Broadcaster
	config      Config
	mu          sync.RWMutex
}

func (p *PubSub) Clean() ([]Client, *ChangeLog) {
	var clients []string
	now := time.Now().UnixNano()
	p.mu.Lock()
	for clientID, timestamp := range p.inactiveClients {
		if now-timestamp >= int64(p.config.InvalidationTime) {
			clients = append(clients, clientID)
		}
	}
	p.mu.Unlock()
	return p.Disconnect(DisconnectOptions{
		Clients:   clients,
		TimeStamp: time.Now().UnixNano(),
	})
}

func (p *PubSub) Connect(opts ConnectOptions) ConnectResult {
	opts.Validate()
	changelog := newChangeLog(opts.TimeStamp)

	mutator := p.broadcaster.Mutator(opts.TimeStamp)
	defer mutator.Close()

	isReconnected := false

	p.mu.Lock()
	c, clientExists := p.clients[opts.ClientID]

	if clientExists {
		isReconnected = true
		if c.User() != opts.UserID {
			p.mu.Unlock()
			return ConnectResult{
				Client: nil,
				ChangeLog: changelog,
				Reconnected: isReconnected,
				Error: errors.New("failed to connect: user id mismatch"),
			}
		}
		delete(p.inactiveClients, c.ID())
	} else {
		c = client.New(opts.ClientID, opts.UserID)
		p.clients[c.ID()] = c

		changelog.addCreatedClient(c.ID())

		if opts.UserID != "" {
			u, userExists := p.users[opts.UserID]
			if !userExists {
				u = user.New()
				p.users[opts.UserID] = u
				changelog.addCreatedUser(opts.UserID)
			}
			u.Add(c)
			mutator.AttachUser(c.ID(), opts.UserID)
			for topicId, sub := range u.Subscriptions().Map() {
				if !sub.Active {
					continue
				}
				// if sub is not active try to add for information about previous user subscriptions
				t := p.topics[topicId]
				if err := c.Subscriptions().Add(topicId, opts.TimeStamp); err == nil {
					t.Add(c.ID())
					mutator.Subscribe(c.ID(), topicId)
				}
			}
		}
	}

	mutator.UpdateClient(c.ID(), opts.Writer)
	p.mu.Unlock()

	return ConnectResult{
		Client: c,
		ChangeLog: changelog,
		Reconnected: isReconnected,
		Error: nil,
	}
}

func (p *PubSub) Inactivate(id string, ts int64) (Client, *ChangeLog, error) {

	changelog := newChangeLog(ts)

	mutator := p.broadcaster.Mutator(ts)
	defer mutator.Close()

	p.mu.Lock()
	defer p.mu.Unlock()

	c, ok := p.clients[id]
	if !ok {
		return nil, changelog, errors.New("failed to inactivate: no client with such id")
	}
	p.inactiveClients[id] = ts

	changelog.addInactivatedClient(id)

	mutator.UpdateClient(id, nil)

	return c, changelog, nil
}

func (p *PubSub) Disconnect(opts DisconnectOptions) ([]Client, *ChangeLog) {
	opts.Validate()
	changelog := newChangeLog(opts.TimeStamp)

	var clients []Client

	mutator := p.broadcaster.Mutator(opts.TimeStamp)
	defer mutator.Close()

	p.mu.Lock()
	defer p.mu.Unlock()
	if opts.All {
		for clientID, c := range p.clients {
			mutator.DeleteClient(clientID)
			changelog.addDeletedClient(clientID)
			clients = append(clients, c)
			c.Invalidate()
			delete(p.clients, clientID)
			delete(p.inactiveClients, clientID)

			if c.User() != "" {
				if u, userExists := p.users[c.User()]; userExists {
					u.Del(c.ID())
					mutator.DetachUser(c.ID(), c.User(), true)
					if len(u.Clients()) == 0 {
						changelog.addDeletedUser(c.User())
						delete(p.users, c.User())
					}
				}
			}
			for topicID := range c.Subscriptions().Map() {
				mutator.Unsubscribe(c.ID(), topicID, true)
				if t, topicExists := p.topics[topicID]; topicExists {
					t.Del(topicID)
					if t.Len() == 0 {
						changelog.addDeletedTopic(topicID)
						delete(p.topics, topicID)
					}
				}
			}
		}
	} else {
		for _, clientID := range opts.Clients {
			mutator.DeleteClient(clientID)
			c, clientExists := p.clients[clientID]
			delete(p.inactiveClients, clientID)
			if !clientExists {
				continue
			}

			changelog.addDeletedClient(clientID)
			clients = append(clients, c)
			c.Invalidate()
			delete(p.clients, clientID)
			if c.User() != "" {
				if u, userExists := p.users[c.User()]; userExists {
					u.Del(c.ID())
					mutator.DetachUser(c.ID(), c.User(), true)
					if len(u.Clients()) == 0 {
						changelog.addDeletedUser(c.User())
						delete(p.users, c.User())
					}
				}
			}
			for topicID := range c.Subscriptions().Map() {
				mutator.Unsubscribe(c.ID(), topicID, true)
				if t, topicExists := p.topics[topicID]; topicExists {
					t.Del(c.ID())
					if t.Len() == 0 {
						changelog.addDeletedTopic(topicID)
						delete(p.topics, topicID)
					}
				}
			}
		}
		for _, userID := range opts.Users {
			u, userExists := p.users[userID]

			if !userExists {
				continue
			}

			delete(p.users, userID)
			changelog.addDeletedUser(userID)

			for clientID, c := range u.Clients() {
				u.Del(clientID)

				mutator.DetachUser(c.ID(), userID, true)

				c.Invalidate()

				delete(p.clients, clientID)

				changelog.addDeletedClient(c.ID())

				clients = append(clients, c)

				mutator.DeleteClient(c.ID())

				for topicID := range c.Subscriptions().Map() {
					mutator.Unsubscribe(c.ID(), topicID, true)
					if t, topicExists := p.topics[topicID]; topicExists {
						t.Del(topicID)
						if t.Len() == 0 {
							changelog.addDeletedTopic(topicID)
							delete(p.topics, topicID)
						}
					}
				}
			}
		}
	}
	return clients, changelog
}

func (p *PubSub) SubscribeClient(opts SubscribeClientOptions) (*ChangeLog, error) {
	opts.Validate()
	changelog := newChangeLog(opts.TimeStamp)

	mut := p.broadcaster.Mutator(opts.TimeStamp)
	defer mut.Close()

	p.mu.Lock()
	defer p.mu.Unlock()

	topicID := opts.Topic

	t, topicExists := p.topics[topicID]

	if !topicExists {
		t = topic.New()
	}

	c, clientExists := p.clients[opts.ID]

	if !clientExists {
		return changelog, errors.New("failed to find a client with specified id")
	}

	if err := c.Subscriptions().Add(topicID, opts.TimeStamp); err == nil {
		t.Add(c.ID())
		mut.Subscribe(c.ID(), topicID)
	} else {
		return changelog, err
	}

	if !topicExists && t.Len() > 0 {
		changelog.addCreatedTopic(topicID)
		p.topics[topicID] = t
	}

	return changelog, nil
}

func (p *PubSub) SubscribeUser(opts SubscribeUserOptions) (*ChangeLog, error) {
	opts.Validate()
	changelog := newChangeLog(opts.TimeStamp)

	mut := p.broadcaster.Mutator(opts.TimeStamp)
	defer mut.Close()

	p.mu.Lock()
	defer p.mu.Unlock()

	topicID := opts.Topic

	t, topicExists := p.topics[topicID]

	if !topicExists {
		t = topic.New()
	}

	u, userExists := p.users[opts.ID]
	if !userExists {
		return changelog, errors.New("failed to find a user with specified id")
	}
	err := u.Subscribe(topicID, opts.TimeStamp, func(c *client.Client) {
		t.Add(c.ID())
		mut.Subscribe(c.ID(), topicID)
	})

	return changelog, err
}

func (p *PubSub) UnsubscribeClient(opts UnsubscribeClientOptions) (*ChangeLog, error) {
	opts.Validate()
	changelog := newChangeLog(opts.TimeStamp)
	mutator := p.broadcaster.Mutator(opts.TimeStamp)
	defer mutator.Close()
	p.mu.Lock()
	defer p.mu.Unlock()
	u, ok := p.users[opts.ID]
	if !ok {
		return changelog, errors.New("failed to find a user with a specified id")
	}
	if opts.All {
		for topicId := range u.Subscriptions().Map() {
			u.Unsubscribe(topicId, opts.TimeStamp, func(c *client.Client) {
				if t, ok := p.topics[topicId]; ok {
					t.Del(c.ID())
					if t.Len() == 0 {
						delete(p.topics, topicId)
						changelog.addDeletedTopic(topicId)
					}
				}
				mutator.Unsubscribe(c.ID(), topicId, false)
			})
		}
		return changelog, nil
	} else {
		err := u.Unsubscribe(opts.Topic, opts.TimeStamp, func(c *client.Client) {
			t, ok := p.topics[opts.Topic]
			if ok {
				t.Del(c.ID())
				mutator.Unsubscribe(c.ID(), opts.Topic, false)
				if t.Len() == 0 {
					delete(p.topics, opts.Topic)
					changelog.addDeletedTopic(opts.Topic)
				}
			}
		})
		return changelog, err
	}
}

func (p *PubSub) UnsubscribeUser(opts UnsubscribeUserOptions) (*ChangeLog, error) {
	opts.Validate()
	changelog := newChangeLog(opts.TimeStamp)
	mutator := p.broadcaster.Mutator(opts.TimeStamp)
	defer mutator.Close()
	p.mu.Lock()
	defer p.mu.Unlock()
	u, ok := p.users[opts.ID]
	if !ok {
		return changelog, errors.New("failed to find client with a specified id")
	}
	if opts.All {
		for topicId := range u.Subscriptions().Map() {
			u.Unsubscribe(topicId, opts.TimeStamp, func(c *client.Client) {
				if t, ok := p.topics[topicId]; ok {
					t.Del(c.ID())
					if t.Len() == 0 {
						delete(p.topics, topicId)
						changelog.addDeletedTopic(topicId)
					}
				}
				mutator.Unsubscribe(c.ID(), topicId, false)
			})
		}
		return changelog, nil
	} else {
		err := u.Unsubscribe(opts.Topic, opts.TimeStamp, func(c *client.Client) {
			if t, ok := p.topics[opts.Topic]; ok {
				t.Del(c.ID())
				if t.Len() == 0 {
					delete(p.topics, opts.Topic)
					changelog.addDeletedTopic(opts.Topic)
				}
			}
			mutator.Unsubscribe(c.ID(), opts.Topic, false)
		})
		return changelog, err
	}
}

func (p *PubSub) Publish(opts PublishOptions) {
	opts.Validate()
	p.broadcaster.Broadcast(opts.Clients, opts.Users, opts.Topics, opts.Data)
}

func (p *PubSub) Start(ctx context.Context) {
	p.broadcaster.Start(ctx)
}

func (p *PubSub) Metrics() Metrics {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return Metrics{
		Clients:         len(p.clients),
		Users:           len(p.users),
		Topics:          len(p.topics),
		InactiveClients: len(p.inactiveClients),
	}
}

func (p *PubSub) ClientSubscribed(client string, topic string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if cl, ok := p.clients[client]; ok {
		sub, ok := cl.Subscriptions().Map()[topic]
		return sub.Active && ok, nil
	}
	return false, errors.New("no client found")
}

func (p *PubSub) UserSubscribed(user string, topic string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if u, ok := p.users[user]; ok {
		sub, ok := u.Subscriptions().Map()[topic]
		return sub.Active && ok, nil
	}
	return false, errors.New("no client found")
}

func (p *PubSub) TopicSubscribers(topic string) ([]string, error) {
	var subs []string
	p.mu.RLock()
	defer p.mu.RUnlock()
	if t, ok := p.topics[topic]; ok {
		subs = make([]string, 0, t.Len())
		for id, _ := range t.Clients() {
			subs = append(subs, id)
		}
		return subs, nil
	}

	return nil, errors.New("no topic found")

}

func (p *PubSub) ClientSubscriptions(id string) ([]string, error) {
	var subs []string
	p.mu.RLock()
	defer p.mu.RUnlock()
	if t, ok := p.clients[id]; ok {
		subs = make([]string, 0, len(t.Subscriptions().Map()))
		for topicID, sub := range t.Subscriptions().Map() {
			if sub.Active {
				subs = append(subs, topicID)
			}
		}
		return subs, nil
	}
	return nil, errors.New("no client found")
}

func (p *PubSub) UserSubscriptions(id string) ([]string, error) {
	var subs []string
	p.mu.RLock()
	defer p.mu.RUnlock()
	if t, ok := p.users[id]; ok {
		subs = make([]string, 0, len(t.Subscriptions().Map()))
		for topicID, sub := range t.Subscriptions().Map() {
			if sub.Active {
				subs = append(subs, topicID)
			}
		}

		return subs, nil
	}
	return nil, errors.New("no user found")
}

func (p *PubSub) UserClients(id string) ([]string, error) {
	var clients []string
	p.mu.RLock()
	defer p.mu.RUnlock()
	if t, ok := p.users[id]; ok {
		clients = make([]string, 0, len(t.Clients()))
		for clientID := range t.Clients() {
			clients = append(clients, clientID)
		}
		return clients, nil
	}
	return nil, errors.New("no user found")
}

func (p *PubSub) Client(id string) (Client, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	c, ok := p.clients[id]
	if ok {
		return c, nil
	}
	return nil, errors.New("no client found")
}

func New(cfg Config) *PubSub {
	if cfg.ProtoProvider == nil {
		cfg.ProtoProvider = batchproto.NewProvider(1024)
	}

	return &PubSub{
		clients:         map[string]*client.Client{},
		users:           map[string]*user.User{},
		topics:          map[string]*topic.Topic{},
		inactiveClients: map[string]int64{},
		config:          cfg,
		broadcaster:     broadcaster.New(cfg.ProtoProvider),
	}
}
