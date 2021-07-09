package engine

import (
	"errors"
	"github.com/ischenkx/swirl/internal/pubsub/common"
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"github.com/ischenkx/swirl/internal/pubsub/session"
	"github.com/ischenkx/swirl/internal/pubsub/subscription"
	"io"
	"sync"
	"time"
)

type Engine struct {
	clients map[string]*client
	topics  map[string]*topic
	presenceManager *presenceManager

	gc *garbageCollector

	queue chan common.Flusher

	mu sync.RWMutex
}

func (e *Engine) enqueue(d common.Flusher) {
	e.queue <- d
}

func (e *Engine) Clean() []ClientInfo {
	e.mu.Lock()
	defer e.mu.Unlock()
	infoAggregator := make([]ClientInfo, 0)
	e.gc.collect(func(clientID string) {
		e.presenceManager.delete(clientID)
		c, ok := e.clients[clientID]
		if !ok {
			return
		}
		delete(e.clients, clientID)

		infoAggregator = append(infoAggregator, ClientInfo{
			Subscriptions: c.subscriptions,
			ID:            clientID,
		})
	})
	return infoAggregator
}

func (e *Engine) IsActive(id string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.presenceManager.isActive(id)
}

func (e *Engine) Inactivate(id string, ts int64) {
	e.mu.Lock()
	c, ok := e.clients[id]
	if !ok {
		e.mu.Unlock()
		return
	}
	e.presenceManager.update(id, false, ts)
	e.gc.add(id, ts)
	e.mu.Unlock()
	c.session.Update(nil, ts, false)
}

func (e *Engine) Disconnect(id string) (ClientInfo, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	c, ok := e.clients[id]
	if !ok {
		return ClientInfo{}, false
	}

	e.presenceManager.delete(id)
	e.gc.delete(id, time.Now().UnixNano())
	delete(e.clients, id)

	c.session.Close()
	return ClientInfo{
		Subscriptions: c.subscriptions,
		ID:            id,
	}, true
}

func (e *Engine) SendToClient(id string, messages []message.Message) error {
	e.mu.RLock()
	c, ok := e.clients[id]
	e.mu.RUnlock()

	if !ok {
		return errors.New("failed to find such client")
	}

	if c.session.Push(messages, true) {
		e.enqueue(c.session)
	}

	return nil
}

func (e *Engine) SendToTopic(id string, mes message.Message) {
	e.mu.RLock()
	u, ok := e.topics[id]
	e.mu.RUnlock()
	if ok {
		if u.group.Push([]message.Message{mes}, true) {
			e.enqueue(u.group)
		}
	}
}

func (e *Engine) AddClientSubscription(id string, topic string, ts int64) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	c, ok := e.clients[id]
	if !ok {
		return errors.New("no such client")
	}
	if !c.subscriptions.Add(topic, ts) {
		return errors.New("failed to add subscription")
	}
	return nil
}

func (e *Engine) Session(id string) (*session.Session, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	c, ok := e.clients[id]
	if !ok {
		return nil, false
	}
	return c.session, true
}

func (e *Engine) DeleteClientSubscription(id string, topic string, ts int64) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	c, ok := e.clients[id]

	if !ok {
		return errors.New("failed to find such client")
	}

	if !c.subscriptions.Delete(topic, ts) {
		return errors.New("no subscription found")
	}

	return nil
}

func (e *Engine) UpdateClient(id string, w io.WriteCloser, ts int64) (*session.Session, bool, int64, error) {
	var disconnectTime int64
	e.mu.Lock()
	c, clientExists := e.clients[id]
	if !clientExists {
		if w == nil {
			e.mu.Unlock()
			return nil, false, 0, errors.New("no such client and writer is nil")
		}
		c = newClient(id, ts)
		e.clients[id] = c
	} else {
		disconnectTime = e.presenceManager.lastTouch(id)
	}
	e.presenceManager.update(id, w != nil, ts)
	if w == nil {
		e.gc.add(id, ts)
	} else {
		e.gc.delete(id, ts)
	}

	e.mu.Unlock()
	if c.session.Update(w, ts, true) {
		e.enqueue(c.session)
	}
	return c.session, clientExists, disconnectTime, nil
}

func (e *Engine) AddTopicClient(topicID, clientID string, sess *session.Session, ts int64) (topicCreated bool) {
	e.mu.Lock()
	t, ok := e.topics[topicID]
	if !ok {
		t = newTopic()
		e.topics[topicID] = t
	}

	prevLen := len(t.clients)

	t.addClient(clientID, ts)
	if len(t.clients) == 1 && prevLen == 0 {
		topicCreated = true
	}
	e.mu.Unlock()
	t.group.Add(sess, ts)
	return
}

func (e *Engine) DeleteTopicClient(topicID, clientID string, forced bool, ts int64) (bool, int) {
	e.mu.Lock()
	t, ok := e.topics[topicID]
	if !ok {
		e.mu.Unlock()
		return false, 0
	}
	t.deleteClient(clientID, forced, ts)
	clientsAmount := len(t.clients)
	e.mu.Unlock()
	t.group.Delete(clientID, forced, ts)
	return true, clientsAmount
}

func (e *Engine) TopicSubscribers(id string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	t, ok := e.topics[id]

	if !ok {
		return nil
	}

	ids := make([]string, 0, len(t.clients))

	for clientID, info := range t.clients {
		if info.Ok {
			ids = append(ids, clientID)
		}
	}
	return ids
}

func (e *Engine) CountTopicSubscribers(id string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	t, ok := e.topics[id]
	if !ok {
		return 0
	}
	return len(t.clients)
}

func (e *Engine) CountSubscriptions(id string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	c, ok := e.clients[id]
	if !ok {
		return 0
	}
	return c.subscriptions.Count()
}

func (e *Engine) Subscriptions(id string) []subscription.Subscription {
	e.mu.RLock()
	defer e.mu.RUnlock()

	c, ok := e.clients[id]

	if !ok {
		return nil
	}

	ids := make([]subscription.Subscription, 0, c.subscriptions.Count())

	c.subscriptions.Iter(func(s subscription.Subscription) {
		ids = append(ids, s)
	})

	return ids
}

func New(q chan common.Flusher, clientInvTime time.Duration) *Engine {
	return &Engine{
		clients:         map[string]*client{},
		topics:          map[string]*topic{},
		presenceManager: newPresenceManager(),
		gc:              newGarbageCollector(int64(clientInvTime)),
		queue:           q,
	}
}
