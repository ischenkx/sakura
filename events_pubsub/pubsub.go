package eventsps

import (
	"sync"
)
//todo
//make listener non-blocking
//CHAN - listening channel
//BUFFER_SIZE - size of channel buffer
//if length of CHAN equals BUFFER_SIZE then create a proxy channel,
//store it somewhere, start new goroutine and read from proxy to
//CHAN in it. Later, if length of CHAN is less or equal than a half of BUFFER_SIZE
//then we can close proxy and write data directly to CHAN

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