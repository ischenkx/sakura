package eventps

import (
	"sync"
)

const BufferSize = 512

type Pubsub struct {
	subscribers map[chan Event]struct{}
	mu sync.RWMutex
}

func (pubsub *Pubsub) Subscribe() Subscription {
	ch := make(chan Event, BufferSize)
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	pubsub.subscribers[ch] = struct{}{}
	return Subscription{pubsub, ch}
}

func (pubsub *Pubsub) Publish(e Event) {
	pubsub.mu.RLock()
	defer pubsub.mu.RUnlock()
	for sub := range pubsub.subscribers {
		sub <- e
	}
}

func (pubsub *Pubsub) deleteSub(c chan Event) {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	close(c)
	delete(pubsub.subscribers, c)
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		subscribers: map[chan Event]struct{}{},
	}
}