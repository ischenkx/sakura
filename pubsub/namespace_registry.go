package pubsub

import (
	"strings"
	"sync"
)

const (
	Any = 0
)

var DefaultNamespaceConfig = NamespaceConfig{}

type NamespaceConfig struct {
	MaxClients, MaxUsers int
}

type namespaceRegistry struct {
	namespaces map[string]NamespaceConfig
	mu sync.RWMutex
}

func (r *namespaceRegistry) generate(topicID string) topic {
	_, ns := parseTopicData(topicID)
	cfg := DefaultNamespaceConfig
	if ns != "" {
		if rcfg, ok := r.Get(ns); ok {
			cfg = rcfg
		}
	}
	return newTopic(cfg)
}

func (r *namespaceRegistry) Register(ns string, cfg NamespaceConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.namespaces[ns] = cfg
}

func (r *namespaceRegistry) Unregister(ns string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.namespaces, ns)
}

func (r *namespaceRegistry) Get(ns string) (NamespaceConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg, ok := r.namespaces[ns]
	return cfg, ok
}

func newNamespaceRegistry() *namespaceRegistry {
	return &namespaceRegistry{
		namespaces: map[string]NamespaceConfig{},
	}
}

func parseTopicData(topic string) (id, ns string) {
	i := strings.Index(topic, ":")

	if i < 0 {
		return topic, ""
	}

	return topic, topic[:i]
}