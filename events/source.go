package events

import (
	"github.com/google/uuid"
	"sync"
)


//todo
//implement non-blocking event source (probably a goroutine for each handler)
type Source struct {
	handlers map[string]Handler

	mu sync.RWMutex
}

func (s *Source) Emit(e Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, handler := range s.handlers {
		select {
		case handler.events <- e:
		}
	}
}

func (s *Source) MultiHandle(h func(event Event), grs int) HandlerCloser {
	if h == nil {
		return HandlerCloser{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	uid := uuid.New().String()
	closer := HandlerCloser{uid, s}
	closerCh := make(chan struct{})
	eventsCh := make(chan Event, 4096)
	s.handlers[uid] = Handler{
		closer: closerCh,
		events: eventsCh,
	}
	for i := 0; i < grs; i++ {
		go func(h func(Event), closer chan struct{}, events chan Event) {
			for {
				select {
				case <-closer:
					return
				case event := <-events:
					h(event)
				}
			}
		}(h, closerCh, eventsCh)
	}

	return closer
}

func (s *Source) Handle(h func(Event)) HandlerCloser {
	if h == nil {
		return HandlerCloser{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	uid := uuid.New().String()
	closer := HandlerCloser{uid, s}
	closerCh := make(chan struct{})
	eventsCh := make(chan Event, 4096)
	s.handlers[uid] = Handler{
		closer: closerCh,
		events: eventsCh,
	}
	go func(h func(Event), closer chan struct{}, events chan Event) {
		for {
			select {
			case <-closer:
				return
			case event := <-events:
				h(event)
			}
		}
	}(h, closerCh, eventsCh)

	return closer
}

func NewSource() *Source {
	return &Source{
		handlers: map[string]Handler{},
	}
}