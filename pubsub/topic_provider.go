package pubsub

import (
	"hash/fnv"
	"sync"
)

type topicProvider struct {
	locks []*sync.RWMutex
	buckets []map[string]map[int]struct{}
}

func (p *topicProvider) hash(d []byte) (int, error) {
	h := fnv.New32a()
	if _, err := h.Write(d); err != nil {
		return 0, err
	}
	return int(h.Sum32()), nil
}

func (p *topicProvider) add(topics []string, s int) {
	for _, topic := range topics {
		if h, err := p.hash([]byte(topic)); err == nil {
			idx := h % len(p.buckets)
			b := p.buckets[idx]
			mu := p.locks[idx]
			mu.Lock()
			t, ok := b[topic]
			if !ok {
				t = map[int]struct{}{}
				b[topic] = t
			}
			t[s] = struct{}{}
			mu.Unlock()
		}
	}
}

func (p *topicProvider) del(topics []string, s int) {
	for _, topic := range topics {
		if h, err := p.hash([]byte(topic)); err == nil {
			idx := h % len(p.buckets)
			b := p.buckets[idx]
			mu := p.locks[idx]
			mu.Lock()
			delete(b[topic], s)
			mu.Unlock()
		}
	}
}

func (p *topicProvider) get(topics []string) map[int][]string {
	m := map[int][]string{}
	for _, topic := range topics {
		if h, err := p.hash([]byte(topic)); err == nil {
			idx := h % len(p.buckets)
			b := p.buckets[idx]
			mu := p.locks[idx]
			mu.RLock()
			t := b[topic]
			for s := range t {
				m[s] = append(m[s], topic)
			}
			mu.RUnlock()
		}
	}
	return m
}

func newTopicProvider(bs int) topicProvider {
	buckets := make([]map[string]map[int]struct{}, bs)

	for i := range buckets {
		buckets[i] = map[string]map[int]struct{}{}
	}

	locks := make([]*sync.RWMutex, bs)

	for i := range locks {
		locks[i] = &sync.RWMutex{}
	}

	return topicProvider{
		locks:   locks,
		buckets: buckets,
	}
}