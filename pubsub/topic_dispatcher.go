package pubsub

import (
	"hash/fnv"
	"sync"
)

type topicDispatcher struct {
	locks []*sync.RWMutex
	buckets []map[string]map[int]struct{}
}

func (p *topicDispatcher) hash(d []byte) (int, error) {
	h := fnv.New32a()
	if _, err := h.Write(d); err != nil {
		return 0, err
	}
	return int(h.Sum32()), nil
}

func (p *topicDispatcher) Amount() int {
	a := 0

	for i := 0; i < len(p.locks); i++ {
		lock := p.locks[i]
		bucket := p.buckets[i]
		lock.RLock()
		a += len(bucket)
		lock.RUnlock()
	}

	return a
}

func (p *topicDispatcher) processShardResults(c shardResultCollector) (topicsUp, topicsDown []string){
	for s, res := range c.Results() {
		for _, topic := range res.topicsUp {
			if h, err := p.hash([]byte(topic)); err == nil {
				idx := h % len(p.buckets)
				b := p.buckets[idx]
				mu := p.locks[idx]
				mu.Lock()
				t, ok := b[topic]
				if !ok {
					t = map[int]struct{}{}
					b[topic] = t
					topicsUp = append(topicsUp, topic)
				}
				t[s] = struct{}{}
				mu.Unlock()
			}
		}

		for _, topic := range res.topicsDown {
			if h, err := p.hash([]byte(topic)); err == nil {
				idx := h % len(p.buckets)
				b := p.buckets[idx]
				mu := p.locks[idx]
				mu.Lock()
				t := b[topic]
				if len(t) > 0 {
					delete(t, s)
					if len(t) == 0 {
						delete(b, topic)
						topicsDown = append(topicsDown, topic)
					}
				}
				mu.Unlock()
			}
		}
	}
	return
}

func (p *topicDispatcher) get(topics []string) map[int][]string {
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

func newTopicDispatcher(bs int) topicDispatcher {
	buckets := make([]map[string]map[int]struct{}, bs)

	for i := range buckets {
		buckets[i] = map[string]map[int]struct{}{}
	}

	locks := make([]*sync.RWMutex, bs)

	for i := range locks {
		locks[i] = &sync.RWMutex{}
	}

	return topicDispatcher{
		locks:   locks,
		buckets: buckets,
	}
}