package hooks

import (
	"sync"
)

type Hub struct {
	id       int64
	priority Priority
	registry *Registry
	handlers map[string][]handler
	mu       sync.RWMutex
}

func (h *Hub) Emit(name string, args []interface{}, resultHandler func(error, []interface{}) bool) {
	h.mu.RLock()
	var handlers []handler
	if hs, ok := h.handlers[name]; ok {
		handlers = make([]handler, len(hs))
		copy(handlers, hs)
	}
	h.mu.RUnlock()

	if handlers == nil {
		return
	}
	for _, hnd := range handlers {
		res, err := hnd.Call(args)
		if resultHandler != nil {
			outputArgs := make([]interface{}, len(res))
			for i, re := range res {
				outputArgs[i] = re.Interface()
			}
			if !resultHandler(err, outputArgs) {
				break
			}
		}
	}
}

func (h *Hub) On(ev string, hnd interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	handlers, ok := h.handlers[ev]
	if ok {
		handlers = make([]handler, 0, 1)
	}
	handlers = append(handlers, newHandler(hnd))
	h.handlers[ev] = handlers
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.registry.delete(h.id)
	h.id = -1
	h.handlers = nil
}
