package events

import (
	"container/heap"
	"github.com/RomanIschenko/notify/internal/utils"
	"sync"
)

type sortedHubList []*Hub

func (p *sortedHubList) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *sortedHubList) Push(el interface{}) {
	*p = append(*p, el.(*Hub))
}

func (p *sortedHubList) Pop() interface{} {
	el := (*p)[p.Len()-1]
	*p = (*p)[:p.Len()-1]
	return el
}

func (p *sortedHubList) Len() int {
	return len(*p)
}

func (p *sortedHubList) Less(i, j int) bool {
	return (*p)[i].p < (*p)[j].p
}

type Source struct {
	hubs sortedHubList
	seq int64
	mu sync.RWMutex
}

func (e *Source) delete(id int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, hub := range e.hubs {
		if hub.id == id {
			heap.Remove(&e.hubs, i)
			return
		}
	}
}

func (e *Source) Emit(name string, args ...interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if len(e.hubs) == 0 {
		return
	}
	prevPriority := e.hubs[0].p
	wg := &sync.WaitGroup{}
	for _, hub := range e.hubs {
		if hub.p != prevPriority {
			wg.Wait()
		}
		wg.Add(1)
		go func(h *Hub) {
			h.Emit(name, args)
			wg.Done()
		}(hub)
		prevPriority = hub.p
	}
}

func (e *Source) NewHub(p Priority) *Hub {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.seq++
	h := &Hub{
		id:       e.seq,
		emitter:  e,
		p: p,
		handlers: map[string][]utils.Callable{},
	}
	heap.Push(&e.hubs, h)
	return h
}

func New() *Source {
	return &Source{}
}