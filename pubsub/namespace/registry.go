package namespace

import (
	"github.com/RomanIschenko/notify/pubsub/internal/topic"
	"sync"
)

const (
	Any = 0
)

var DefaultConfig = Config{}

type Config struct {
	MaxClients, MaxUsers int
}

type Registry struct {
	namespaces map[string]Config
	mu sync.RWMutex
}

func (r *Registry) generate(topicID string) topic.Topic {
	_, ns := parseTopicData(topicID)
	cfg := DefaultConfig
	if ns != "" {
		if rcfg, ok := r.Get(ns); ok {
			cfg = rcfg
		}
	}
	return topic.New(cfg)
}

func (r *Registry) Register(ns string, cfg Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.namespaces[ns] = cfg
}

func (r *Registry) Unregister(ns string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.namespaces, ns)
}

func (r *Registry) Get(ns string) (Config, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg, ok := r.namespaces[ns]
	return cfg, ok
}

func NewRegistry() *Registry {
	return &Registry{
		namespaces: map[string]Config{},
	}
}