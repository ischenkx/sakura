package events

import "sync"
//todo
//make listener non-blocking
//CHAN - listening channel
//BUFFER_SIZE - size of channel buffer
//if length of CHAN equals BUFFER_SIZE then create a proxy channel,
//store it somewhere, start new goroutine and read from proxy to
//CHAN in it. Later, if length of CHAN is less or equal than a half of BUFFER_SIZE
//then we can close proxy and write data directly to CHAN

const BufferSize = 512

type Listener struct {
	listeners map[chan Event]struct{}
	mu sync.RWMutex
}

func (l *Listener) Listen() chan Event {
	ch := make(chan Event, BufferSize)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.listeners[ch] = struct{}{}
	return ch
}

func (l *Listener) Close(c chan Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.listeners, c)
}

func (l *Listener) Emit(e Event) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for listener := range l.listeners {
		listener <- e
	}
}

func NewListener() *Listener {
	return &Listener{
		listeners: map[chan Event]struct{}{},
	}
}