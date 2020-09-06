package events

import (
	"github.com/google/uuid"
	"sync"
)

type Source struct {
	handlers map[string]Handler
	mu sync.RWMutex
}

func (s *Source) Emit(e Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, handler := range s.handlers {
		handler(e)
	}
}

func (s *Source) Handle(h Handler) HandlerCloser {
	if h == nil {
		return HandlerCloser{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	uid := uuid.New().String()
	closer := HandlerCloser{uid, s}
	s.handlers[uid] = h
	return closer
}

func NewSource() *Source {
	return &Source{
		handlers: map[string]Handler{},
	}
}