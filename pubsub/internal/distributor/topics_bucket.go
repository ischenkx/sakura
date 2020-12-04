package distributor

import (
	"sync"
)

type topicsBucket struct {
	b map[string]*int2int
	mu sync.RWMutex
}

func (b *topicsBucket) len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.b)
}

func (b *topicsBucket) store(t string, s int, tx int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	topicData, ok := b.b[t]

	if !ok {
		topicData = newInt2Int(10, 0.2)
		b.b[t] = topicData
	}

	time, ok := topicData.Get(int64(s))

	if time > tx {
		return
	}

	topicData.Put(int64(s), tx)
}

func (b *topicsBucket) delete(t string, s int, tx int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	topicData, ok := b.b[t]

	if !ok {
		return
	}

	time, ok := topicData.Get(int64(s))

	if time > tx {
		return
	}

	topicData.Del(int64(s))
	if topicData.Size() == 0 {
		delete(b.b, t)
	}
	// probably can lead to inconsistency
	// but it's such a rare that we can leave it like that
}

func newTopicsBucket() *topicsBucket {
	return &topicsBucket{
		b: map[string]*int2int{},
		mu: sync.RWMutex{},
	}
}

