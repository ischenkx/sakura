package hooks

import (
	"container/heap"
	"sync"
)

type Registry struct {
	hubs sortedHubList
	seq  int64
	mu   sync.RWMutex
}

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
	return (*p)[i].priority < (*p)[j].priority
}

func (e *Registry) delete(id int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, hub := range e.hubs {
		if hub.id == id {
			heap.Remove(&e.hubs, i)
			return
		}
	}
}

func (e *Registry) Emit(name string, args []interface{}, resultHandler func(error, []interface{}) bool) {
	e.mu.RLock()
	if len(e.hubs) == 0 {
		e.mu.RUnlock()
		return
	}
	hubs := make(sortedHubList, len(e.hubs))
	copy(hubs, e.hubs)
	e.mu.RUnlock()
	for _, hub := range e.hubs {
		hub.Emit(name, args, resultHandler)
	}
}

func (e *Registry) NewHub(p Priority) *Hub {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.seq++
	h := &Hub{
		id:       e.seq,
		registry: e,
		priority: p,
		handlers: map[string][]handler{},
	}
	heap.Push(&e.hubs, h)
	return h
}

func NewRegistry() *Registry {
	return &Registry{}
}
