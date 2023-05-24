package broadcaster

import "sync"

type SubscriptionManager struct {
	topics map[string]map[string]struct{}
	users  map[string]map[string]struct{}
	mu     sync.RWMutex
}

func newSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		topics: map[string]map[string]struct{}{},
		users:  map[string]map[string]struct{}{},
		mu:     sync.RWMutex{},
	}
}

func (manager *SubscriptionManager) Add(topic, user string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.initUser(user)
	manager.initTopic(topic)

	manager.topics[topic][user] = struct{}{}
	manager.users[user][topic] = struct{}{}
}

func (manager *SubscriptionManager) Remove(topic, user string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if topicInfo, ok := manager.topics[topic]; ok {
		delete(topicInfo, topic)
	}

	if userInfo, ok := manager.users[user]; ok {
		delete(userInfo, user)
	}
}

func (manager *SubscriptionManager) Iter(topic string, iter func(string)) {
	var ids []string

	manager.mu.RLock()
	if data, ok := manager.topics[topic]; ok {
		for id := range data {
			ids = append(ids, id)
		}
	}
	manager.mu.RUnlock()

	for _, id := range ids {
		iter(id)
	}
}

func (manager *SubscriptionManager) RemoveByUser(id string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if user, ok := manager.users[id]; ok {
		delete(manager.users, id)
		for topicID := range user {
			delete(manager.topics[topicID], id)
		}
	}
}

func (manager *SubscriptionManager) RemoveByTopic(id string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if topic, ok := manager.topics[id]; ok {
		delete(manager.topics, id)
		for userID := range topic {
			delete(manager.users[userID], id)
		}
	}
}

func (manager *SubscriptionManager) initUser(id string) {
	if _, ok := manager.users[id]; !ok {
		manager.users[id] = map[string]struct{}{}
	}
}

func (manager *SubscriptionManager) initTopic(id string) {
	if _, ok := manager.topics[id]; !ok {
		manager.topics[id] = map[string]struct{}{}
	}
}
