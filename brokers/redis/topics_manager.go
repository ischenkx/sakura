package redibroker

import (
	"sync"
	"time"
)

type topicsManager struct {
	topics map[string]struct{}
	inactiveTopics map[string]time.Time
	mu sync.RWMutex
	topicTTL time.Duration
}

func (m *topicsManager) add(ts ...string) []string {
	added := make([]string, 0, len(ts)/4)
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range ts {
		if _, ok := m.topics[t]; ok {
			m.topics[t] = struct{}{}
			delete(m.inactiveTopics, t)
			added = append(added, t)
		} else {
			m.topics[t] = struct{}{}
			added = append(added, t)
		}
	}

	return added
}

func (m *topicsManager) del(ts ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range ts {
		if _, ok := m.topics[t]; ok {
			m.inactiveTopics[t] = time.Now()
		}
	}
}

func (m *topicsManager) clean() (deleted []string) {
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	for topic, t := range m.inactiveTopics {
		if now.Sub(t) >= m.topicTTL {
			delete(m.inactiveTopics, topic)
			delete(m.topics, topic)
			deleted = append(deleted, topic)
		}
	}
	return
}

func newTopicsManager() *topicsManager {
	return &topicsManager{
		topics: map[string]struct{}{},
		inactiveTopics: map[string]time.Time{},
		mu:             sync.RWMutex{},
		topicTTL:       time.Minute * 2,
	}
}