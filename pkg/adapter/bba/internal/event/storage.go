package event

import (
	"github.com/google/uuid"
	"github.com/ischenkx/swirl"
	"sync"
)

type handle struct {
	sto *CustomStorage
	uid string
	event string
}

func (h handle) Close() {
	h.sto.mu.Lock()
	defer h.sto.mu.Unlock()

	arr, ok := h.sto.handlers[h.event]
	if !ok {
		return
	}

	for i := 0; i < len(arr); i++ {
		if arr[i].uid == h.uid {
			arr[i] = arr[len(arr)-1]
			arr = arr[:len(arr) - 1]
		}
	}

	if len(arr) == 0 {
		delete(h.sto.handlers, h.event)
	} else {
		h.sto.handlers[h.event] = arr
	}
}

type callable struct {
	uid string
}

type CustomStorage struct {
	mu sync.RWMutex
	handlers map[string][]*callable
}

func (sto *CustomStorage) Handle(name string, handler interface{}) (swirl.Handle, error) {
	clb, err := newCallable(handler)
	if err != nil {

	}
	sto.mu.Lock()
	defer sto.mu.Unlock()



	sto.handlers[name] = append(sto.handlers[name], handler)
}
