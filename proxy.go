package notify

import (
	"container/heap"
	"errors"
	"sync"
)

type (
	SubscribeProxyHandler func(*App, *SubscribeOptions)
	UnsubscribeProxyHandler func(*App, *UnsubscribeOptions)
	PublishProxyHandler func(*App, *PublishOptions)
	ConnectProxyHandler func(*App, *ConnectOptions)
	DisconnectProxyHandler func(*App, *DisconnectOptions)
)

type Proxy struct {
	id           int64
	closed       bool
	reg          *proxyStorage

	mu           sync.RWMutex
}



func (hub *Proxy) Close() error {
	hub.mu.Lock()
	if hub.closed {
		hub.mu.Unlock()
		return errors.New("Proxy: already closed")
	}
	hub.closed = true
	reg := hub.reg
	id := hub.id
	hub.reg = nil
	hub.onPublish = nil
	hub.onConnect = nil
	hub.onDisconnect = nil
	hub.mu.Unlock()

	reg.mu.Lock()
	delete(reg.hubs, id)
	reg.mu.Unlock()

	return nil
}

type proxyStorage struct {
	hubs map[int64]*Proxy
	idSeq int64
	priority Priority
	mu sync.RWMutex
}

func (s *proxyStorage) iterEvents(h func(e *Proxy)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ev := range s.hubs {
		h(ev)
	}
}

func (s *proxyStorage) new() *Proxy {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.idSeq++
	id := s.idSeq
	ev := &Proxy{
		id:           id,
		closed:       false,
		reg:          s,
		mu:           sync.RWMutex{},
	}
	s.hubs[id] = ev
	return ev
}

type sortablePS []*proxyStorage

func (p *sortablePS) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *sortablePS) Push(el interface{}) {
	*p = append(*p, el.(*proxyStorage))
}

func (p *sortablePS) Pop() interface{} {
	el := (*p)[p.Len()-1]
	*p = (*p)[:p.Len()-1]
	return el
}

func (p *sortablePS) Len() int {
	return len(*p)
}

func (p *sortablePS) Less(i, j int) bool {
	return (*p)[i].priority < (*p)[j].priority
}

type proxyRegistry struct {
	tables []*proxyStorage
	mu     sync.RWMutex
}

func (r *proxyRegistry) new(p Priority) *Proxy {
	r.mu.Lock()
	var sto *proxyStorage
	for _, s := range r.tables {
		if s.priority == p {
			sto = s
			break
		}
	}
	if sto == nil {
		sto = &proxyStorage{
			hubs: map[int64]*Proxy{},
			idSeq:    0,
			priority: p,
		}
		heap.Push((*sortablePS)(&r.tables), sto)
	}
	r.mu.Unlock()
	return sto.new()
}

func (r *proxyRegistry) emitConnect(a *App, opts *ConnectOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Proxy) {
			e.emitConnect(wg, a, opts)
		})
		wg.Wait()
	}
}

func (r *proxyRegistry) emitDisconnect(a *App, opts *DisconnectOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Proxy) {
			e.emitDisconnect(wg, a, opts)
		})
		wg.Wait()
	}
}

func (r *proxyRegistry) emitSubscribe(a *App, opts *SubscribeOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Proxy) {
			e.emitSubscribe(wg, a, opts)
		})
		wg.Wait()
	}
}

func (r *proxyRegistry) emitUnsubscribe(a *App, opts *UnsubscribeOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Proxy) {
			e.emitUnsubscribe(wg, a, opts)
		})
		wg.Wait()
	}
}

func (r *proxyRegistry) emitPublish(a *App, opts *PublishOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Proxy) {
			e.emitPublish(wg, a, opts)
		})
		wg.Wait()
	}
}

func newProxyRegistry() *proxyRegistry {
	return &proxyRegistry{}
}