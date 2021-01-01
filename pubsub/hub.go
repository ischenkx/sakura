package pubsub

import (
	"context"
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify/pubsub/internal/broadcaster"
	"github.com/RomanIschenko/notify/pubsub/internal/client"
	"github.com/RomanIschenko/notify/pubsub/internal/topic"
	"github.com/RomanIschenko/notify/pubsub/internal/user"
	"sync"
	"time"
)

type Config struct {
	Workers int
	ClientTTL time.Duration
	ClientBufferSize int
	MaxBatchSize int
}

type Hub struct {
	config 			Config
	clients 		map[string]*client.Client
	users 			map[string]*user.User
	topics 			map[string]*topic.Topic
	inactiveClients map[string]int64
	broadcaster 	*broadcaster.Broadcaster
	mu 				sync.RWMutex
}

func (h *Hub) subscribe(opts SubscribeOptions) (log ChangeLog) {
	log.setTime()
	h.mu.Lock()
	for _, topicId := range opts.Topics {
		t, topicExists := h.topics[topicId]
		if !topicExists {
			t = topic.New(h.broadcaster.NewGroup())
			h.topics[topicId] = t
		}
		for _, userID := range opts.Users {
			if u, ok := h.users[userID]; ok {
				u.Subscribe(topicId, opts.Seq, func(c *client.Client) {
					h.broadcaster.Join(c.Descriptor(), t.Descriptor())
				})
			}
		}
		for _, clientID := range opts.Clients {
			if c, ok := h.clients[clientID]; ok {
				if err := c.Subscribe(topicId, opts.Seq); err == nil {
					h.broadcaster.Join(c.Descriptor(), t.Descriptor())
				}
			}
		}
		if !topicExists {
			if !h.broadcaster.CheckGroupValidity(t.Descriptor()) {
				delete(h.topics, topicId)
			} else {
				log.TopicsUp = append(log.TopicsUp, topicId)
			}
		}
	}
	h.mu.Unlock()
	return
}

func (h *Hub) unsubscribe(opts UnsubscribeOptions) (log ChangeLog) {
	log.setTime()
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, topicId := range opts.Topics {
		t, topicExists := h.topics[topicId]
		if !topicExists {
			continue
		}
		for _, userID := range opts.Users {
			if u, ok := h.users[userID]; ok {
				u.Unsubscribe(topicId, opts.Seq, func(c *client.Client) {
					if h.broadcaster.Leave(c.Descriptor(), t.Descriptor()) == 0 {
						delete(h.topics, topicId)
						log.TopicsDown = append(log.TopicsDown, topicId)
					}
				})
			}
		}
		for _, clientID := range opts.Clients {
			if c, ok := h.clients[clientID]; ok {
				if err := c.Unsubscribe(topicId, opts.Seq); err == nil {
					if h.broadcaster.Leave(c.Descriptor(), t.Descriptor()) == 0 {
						delete(h.topics, topicId)
						log.TopicsDown = append(log.TopicsDown, topicId)
					}
				}
			}
		}
	}
	return
}

func (h *Hub) publish(opts PublishOptions) {
	descriptors := make([]uint64, 0, len(opts.Users)+len(opts.Clients)+len(opts.Topics))
	h.mu.RLock()
	for _, id := range opts.Clients {
		if c, ok := h.clients[id]; ok {
			descriptors = append(descriptors, c.Descriptor())
		}
	}

	for _, id := range opts.Users {
		if u, ok := h.users[id]; ok {
			descriptors = append(descriptors, u.Descriptor())
		}
	}

	for _, id := range opts.Topics {
		if t, ok := h.topics[id]; ok {
			descriptors = append(descriptors, t.Descriptor())
		}
	}
	h.mu.RUnlock()
	messages := []broadcaster.Message{
		{
			Data:        opts.Message,
			NoBuffering: opts.NoBuffering,
		},
	}
	for _, desc := range descriptors {
		h.broadcaster.Push(desc, messages)
	}
}

func (h *Hub) connect(opts ConnectOptions) (Client, ChangeLog, error) {
	var log ChangeLog
	log.setTime()
	h.mu.Lock()
	defer h.mu.Unlock()
	c, ok := h.clients[opts.ClientID]
	delete(h.inactiveClients, opts.ClientID)
	if ok {
		if c.User() != opts.UserID {
			return Client{}, log, errors.New("user id mismatch")
		}
		h.broadcaster.UpdateSession(c.Descriptor(), opts.Transport)
	} else {
		c = client.New(h.broadcaster.NewSession(opts.Transport), opts.ClientID, opts.UserID)
		h.clients[opts.ClientID] = c
		log.ClientsUp = append(log.ClientsUp, c.ID())
	}
	if c.User() != "" {
		u, ok := h.users[c.User()]
		if !ok {
			u = user.New(h.broadcaster.NewGroup())
			h.users[c.User()] = u
			log.UsersUp = append(log.UsersUp, c.User())
		}
		u.Clients()[opts.ClientID] = c
	}
	return Client{c}, log, nil
}

func (h *Hub) inactivate(opts InactivateOptions) {
	h.mu.Lock()
	c, ok := h.clients[opts.ClientID]
	if ok {
		if _, inactive := h.inactiveClients[c.ID()]; !inactive {
			h.inactiveClients[c.ID()] = time.Now().UnixNano()
		}
	}
	desc := c.Descriptor()
	h.mu.Unlock()
	if ok {
		h.broadcaster.UpdateSession(desc, nil)
	}
}

func (h *Hub) disconnect(opts DisconnectOptions) (log ChangeLog) {
	log.setTime()
	h.mu.Lock()
	defer h.mu.Unlock()

	if opts.All {
		fmt.Println("i am sorry but i haven't implemented disconnection of all clients.")
		return
	}

	for _, id := range opts.Clients {
		c, ok := h.clients[id]
		if !ok {
			continue
		}
		h.broadcaster.DeleteSession(c.Descriptor())
		log.ClientsDown = append(log.ClientsDown, id)
		for topicId, sub := range c.Subscriptions() {
			if !sub.Active() {
				continue
			}
			if t, ok := h.topics[topicId]; ok {
				if h.broadcaster.Leave(c.Descriptor(), t.Descriptor()) == 0 {
					delete(h.topics, topicId)
					log.TopicsDown = append(log.TopicsDown, topicId)
				}
			}
		}
		if c.User() != "" {
			u, ok := h.users[c.User()]
			if !ok {
				continue
			}
			delete(u.Clients(), id)
			if len(u.Clients()) == 0 {
				log.UsersDown = append(log.UsersDown, c.User())
				delete(h.users, c.User())
			}
			h.broadcaster.Leave(c.Descriptor(), u.Descriptor())
		}
	}

	for _, id := range opts.Users {
		if u, ok := h.users[id]; ok {
			delete(h.users, id)
			log.UsersDown = append(log.UsersDown, id)
			for _, c := range u.Clients() {
				h.broadcaster.DeleteSession(c.Descriptor())
				delete(h.clients, c.ID())
				log.ClientsDown = append(log.ClientsDown, c.ID())
				h.broadcaster.Leave(c.Descriptor(), u.Descriptor())
				for topicID, _ := range c.Subscriptions() {
					if t, ok := h.topics[topicID]; ok {
						if h.broadcaster.Leave(c.Descriptor(), t.Descriptor()) == 0 {
							delete(h.topics, topicID)
							log.TopicsDown = append(log.TopicsDown, topicID)
						}
					}
				}
			}
		}
	}
	return
}

func (h *Hub) Inactivate(opts InactivateOptions) {
	h.inactivate(opts)
}

func (h *Hub) Connect(opts ConnectOptions) (Client, ChangeLog, error) {
	return h.connect(opts)
}

func (h *Hub) Disconnect(opts DisconnectOptions) (log ChangeLog) {
	return h.disconnect(opts)
}

func (h *Hub) Publish(opts PublishOptions) {
	h.publish(opts)
}

func (h *Hub) Subscribe(opts SubscribeOptions) (log ChangeLog) {
	return h.subscribe(opts)
}

func (h *Hub) Unsubscribe(opts UnsubscribeOptions) (log ChangeLog) {
	return h.unsubscribe(opts)
}

func (h *Hub) Clean() (log ChangeLog) {
	log.setTime()
	sessions := make([]uint64, 0)
	h.mu.Lock()
	now := time.Now().UnixNano()
	for id, seq := range h.inactiveClients {
		if now - seq < int64(h.config.ClientTTL) {
			continue
		}
		delete(h.inactiveClients, id)
		if c, ok := h.clients[id]; ok {
			delete(h.clients, id)
			log.ClientsDown = append(log.ClientsDown, id)
			sessions = append(sessions, c.Descriptor())
			if c.User() != "" {
				if u, ok := h.users[c.User()]; ok {
					delete(u.Clients(), c.ID())
					h.broadcaster.Leave(c.Descriptor(), u.Descriptor())
					if len(u.Clients()) == 0 {
						log.UsersDown = append(log.UsersDown, c.User())
						delete(h.users, c.User())
					}
				}
			}

			for topicId, sub := range c.Subscriptions() {
				if !sub.Active() {
					continue
				}
				if t, ok := h.topics[topicId]; ok {
					if h.broadcaster.Leave(c.Descriptor(), t.Descriptor()) == 0 {
						log.TopicsDown = append(log.TopicsDown, topicId)
						delete(h.topics, topicId)
					}
				}
			}
		}
	}
	h.mu.Unlock()
	for _, desc := range sessions {
		h.broadcaster.DeleteSession(desc)
	}
	return
}

func (h *Hub) Metrics() (m Metrics) {
	m.BroadcasterMetrics = h.broadcaster.Metrics()
	h.mu.RLock()
	m.Clients = len(h.clients)
	m.Users = len(h.users)
	m.Topics = len(h.topics)
	m.InactiveClients = len(h.inactiveClients)
	h.mu.RUnlock()
	return
}

func (h *Hub) Start(ctx context.Context) {
	h.broadcaster.Start(ctx)
}

func New(cfg Config) *Hub {
	return &Hub{
		config:          cfg,
		clients: map[string]*client.Client{},
		users: map[string]*user.User{},
		topics: map[string]*topic.Topic{},
		inactiveClients: map[string]int64{},
		broadcaster:     broadcaster.New(broadcaster.Config{
			MaxBatchSize:        cfg.MaxBatchSize,
			ClientBufferSize:    cfg.ClientBufferSize,
			SessionsPerSubGroup: 2<<8,
		}),
	}
}