package registry

import (
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"hash/fnv"

	//"hash/fnv"
)

type topicsRegistry []*topicsBucket

func (p *topicsRegistry) hash(d []byte) (int, error) {
	h := fnv.New32a()
	if _, err := h.Write(d); err != nil {
		return 0, err
	}
	return int(h.Sum32()), nil
}

func (p *topicsRegistry) Amount() int {
	a := 0

	for _, b := range *p {
		a += b.len()
	}

	return a
}

func (p *topicsRegistry) processResult(s int, res changelog.Log) {
	for _, topic := range res.TopicsUp {
		if h, err := p.hash([]byte(topic)); err == nil {
			idx := h % len(*p)
			b := (*p)[idx]
			b.store(topic, s, res.Time)
		}
	}

	for _, topic := range res.TopicsDown {
		if h, err := p.hash([]byte(topic)); err == nil {
			idx := h % len(*p)
			b := (*p)[idx]
			b.delete(topic, s, res.Time)
		}
	}
}

func (p *topicsRegistry) loadTopics(topics []string, s *Sharder) {
	for _, topic := range topics {
		if h, err := p.hash([]byte(topic)); err == nil {
			idx := h % len(*p)
			b := (*p)[idx]
			b.mu.RLock()
			if t, ok := b.b[topic]; ok {
				t.Iter(func(shard, tx int64) {
					s.AddTopics(int(shard), topic)
				})
			}
			b.mu.RUnlock()
		}
	}
}

//func (p *topicsRegistry) iter(topics []string, f func(int, []string)) {
//	m := map[int][]string{}
//	for _, topic := range topics {
//		if h, err := p.hash([]byte(topic)); err == nil {
//			idx := h % len(*p)
//			b := (*p)[idx]
//			b.mu.RLock()
//			t := b.b[topic]
//			for s := range t {
//				m[s] = append(m[s], topic)
//			}
//			b.mu.RUnlock()
//		}
//	}
//
//	for shardIdx, b := range m {
//		f(shardIdx, b)
//	}
//}

func newTopicsRegistry(bs int) topicsRegistry {
	buckets := make([]*topicsBucket, bs)
	for i := range buckets {
		buckets[i] = newTopicsBucket()
	}
	return buckets
}
