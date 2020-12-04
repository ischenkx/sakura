package events

import (
	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
	"io"
	"runtime"
	"sync"
)


type Source struct {
	handlers map[string]map[string]func(interface{})
	pool *ants.Pool
	mu sync.RWMutex
}

func (s *Source) Emit(e Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if handlers, ok := s.handlers[e.Type]; ok {
		for _, handler := range handlers {
			handler(e)
		}
	}
}

func (s *Source) Handle(e string, f func(interface{})) io.Closer {
	if f == nil {
		return HandlerCloser{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	uid := uuid.New().String()
	closer := HandlerCloser{uid,e, s}
	handlers, ok := s.handlers[e]
	if !ok {
		handlers = map[string]func(interface{}){}
		s.handlers[e] = handlers
	}
	handlers[uid] = f
	return closer
}

func NewSource() *Source {
	pool, err := ants.NewPool(runtime.NumCPU())
	if err != nil {
		panic(err)
	}
	return &Source{
		handlers: map[string]map[string]func(interface{}){},
		pool: pool,
	}
}