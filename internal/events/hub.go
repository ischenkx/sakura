package events

import (
	"github.com/ischenkx/notify/internal/utils"
	"sync"
)

type Hub struct {
	id int64
	p Priority
	emitter *Source
	handlers map[string][]utils.Callable
	mu sync.RWMutex
}

func (h *Hub) Emit(name string, args []interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if handlers, ok := h.handlers[name]; ok {
		for _, hnd := range handlers {
			hnd.Call(args)
		}
	}
}

func (h *Hub) On(ev string, hnd interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	handlers, ok := h.handlers[ev]
	if ok {
		handlers = make([]utils.Callable, 1)
	}
	handlers = append(handlers, utils.NewCallable(hnd))
	h.handlers[ev] = handlers
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.id > 0 {
		h.emitter.delete(h.id)
		h.id = -1
		h.handlers = nil
	}
}
