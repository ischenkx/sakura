package user

import (
	"errors"
	"github.com/ischenkx/swirl/internal/pubsub/changelog"
	"github.com/ischenkx/swirl/internal/pubsub/subscription"
	"sync"
)

var subscriptionArrayPool = &sync.Pool{
	New: func() interface{} {
		return make([]subscription.Subscription, 0, 100)
	},
}

var stringsPool = &sync.Pool{
	New: func() interface{} {
		return make([]string, 0, 100)
	},
}

type user struct {
	clients       map[string]struct{}
	subscriptions *subscription.List
}

type Registry struct {
	presenceManager *presenceManager
	users       map[string]user
	client2user map[string]string
	topic2users map[string]map[string]struct{}
	mu          sync.RWMutex
}

func (r *Registry) UserByClient(id string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client2user[id]
}

func (r *Registry) UpdateClientPresence(id string, isActive bool, ts int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	userID, ok := r.client2user[id]
	if !ok {
		return
	}
	r.presenceManager.update(userID, id, isActive, ts)
}

func (r *Registry) Add(userId, clientId string, ts int64) changelog.ChangeLog {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, userExists := r.users[userId]

	if !userExists {
		u = user{
			clients:       map[string]struct{}{},
			subscriptions: subscription.NewList(),
		}
		r.users[userId] = u
	}
	r.client2user[clientId] = userId
	u.clients[clientId] = struct{}{}

	cl := changelog.ChangeLog{
		TimeStamp: ts,
	}

	if !userExists {
		cl.UsersUp = append(cl.UsersUp, userId)
	}

	return cl
}

func (r *Registry) IsActive(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.presenceManager.isActive(id)
}

func (r *Registry) Delete(userId, client string, ts int64) changelog.ChangeLog {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[userId]
	if !ok {
		return changelog.ChangeLog{}
	}

	r.presenceManager.delete(userId, client)
	delete(r.client2user, client)
	delete(u.clients, client)

	cl := changelog.ChangeLog{
		TimeStamp: ts,
	}

	if len(u.clients) == 0 {
		cl.UsersDown = append(cl.UsersDown, userId)
	}

	return cl
}

func (r *Registry) Subscribe(userId, topic string, ts int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[userId]

	if !ok {
		u = user{
			clients:       map[string]struct{}{},
			subscriptions: subscription.NewList(),
		}
		r.users[userId] = u
	}
	if u.subscriptions.Add(topic, ts) {
		t, ok := r.topic2users[topic]
		if !ok {
			t = map[string]struct{}{}
			r.topic2users[topic] = t
		}
		t[userId] = struct{}{}
	} else {
		return errors.New("failed to subcribe")
	}

	return nil
}

func (r *Registry) Unsubscribe(userId, topic string, ts int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[userId]

	if !ok {
		u = user{
			clients:       map[string]struct{}{},
			subscriptions: subscription.NewList(),
		}
		r.users[userId] = u
	}
	if u.subscriptions.Delete(topic, ts) {
		delete(r.topic2users[userId], userId)
	} else {
		return errors.New("failed to unsubscribe")
	}
	return nil
}

func (r *Registry) Iter(id string, f func(id string)) {
	r.mu.RLock()

	u, ok := r.users[id]

	if !ok {
		r.mu.RUnlock()
		return
	}
	ids := stringsPool.Get().([]string)
	for client := range u.clients {
		ids = append(ids, client)
	}
	r.mu.RUnlock()

	for _, client := range ids {
		f(client)
	}

	stringsPool.Put(ids[:0])
}

func (r *Registry) IterSubscriptions(id string, f func(subscription.Subscription)) {
	r.mu.RLock()
	u, ok := r.users[id]
	if !ok {
		r.mu.RUnlock()
		return
	}
	subs := subscriptionArrayPool.Get().([]subscription.Subscription)
	u.subscriptions.Iter(func(s subscription.Subscription) {
		subs = append(subs, s)
	})
	r.mu.RUnlock()

	for _, topic := range subs {
		f(topic)
	}

	subscriptionArrayPool.Put(subs[:0])
}

func (r *Registry) CountSubscriptions(id string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return 0
	}
	return u.subscriptions.Count()
}

func (r *Registry) CountTopicUsers(topic string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.topic2users[topic])
}

func (r *Registry) IterTopicUsers(topic string, f func(id string)) {
	r.mu.RLock()
	users, ok := r.topic2users[topic]
	if !ok {
		r.mu.RUnlock()
		return
	}
	ids := stringsPool.Get().([]string)
	for id := range users {
		ids = append(ids, id)
	}
	r.mu.RUnlock()
	for _, id := range ids {
		f(id)
	}
	stringsPool.Put(ids[:0])
}

func NewRegistry() *Registry {
	return &Registry{
		presenceManager: &presenceManager{users: map[string]map[string]presenceInfo{}},
		users:           map[string]user{},
		client2user:     map[string]string{},
		topic2users:     map[string]map[string]struct{}{},
		mu:              sync.RWMutex{},
	}
}
