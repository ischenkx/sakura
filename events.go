package notify

import (
	"container/heap"
	"errors"
	"sync"
)

// TODO
// Proxy 2.0
// Connect - first connection of the Client
// Reconnect - Client reconnected
// Disconnect - Client disconnected and was deleted
// Inactivate - Client disconnected and is being waited to reconnect
// Message - incoming message
// Change - all changes made in pub/sub system (ChangeLog)
// Publish - new publication
// Subscribe - subscribed
// Unsubscribe - unsubscribed

type (
	PublishHandler    func(a *App, opts PublishOptions)
	ConnectHandler    func(a *App, opts ConnectOptions, c Client)
	ReconnectHandler  func(a *App, opts ConnectOptions, c Client)
	DisconnectHandler func(a *App, cs Client)
	InactivateHandler func(a *App, c Client)
	MessageHandler    func(a *App, c Client, message []byte)
	ChangeHandler     func(a *App, log ChangeLog)
)

type Events struct {
	id           int64
	closed       bool
	reg          *eventsStorage
	onPublish    []PublishHandler
	onConnect    []ConnectHandler
	onReconnect  []ReconnectHandler
	onDisconnect []DisconnectHandler
	onInactivate []InactivateHandler
	onMessage    []MessageHandler
	onChange     []ChangeHandler
	mu           sync.RWMutex
}

func (hub *Events) emitReconnect(wg *sync.WaitGroup, a *App, opts ConnectOptions, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onReconnect {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts, c)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) emitPublish(wg *sync.WaitGroup, a *App, opts PublishOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onPublish {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) emitMessage(wg *sync.WaitGroup, a *App, c Client, data []byte) {
	hub.mu.RLock()
	for _, handler := range hub.onMessage {
		h := handler
		wg.Add(1)
		go func() {
			h(a, c, data)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) emitConnect(wg *sync.WaitGroup, a *App, opts ConnectOptions, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onConnect {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts, c)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) emitDisconnect(wg *sync.WaitGroup, a *App, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onDisconnect {
		h := handler
		wg.Add(1)
		go func() {
			h(a, c)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) emitInactivate(wg *sync.WaitGroup, a *App, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onInactivate {
		h := handler
		wg.Add(1)
		go func() {
			h(a, c)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) emitChange(wg *sync.WaitGroup, a *App, c ChangeLog) {
	hub.mu.RLock()
	for _, handler := range hub.onChange {
		h := handler
		wg.Add(1)
		go func() {
			h(a, c)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Events) OnConnect(f ConnectHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onConnect = append(hub.onConnect, f)
	return nil
}

func (hub *Events) OnPublish(f PublishHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onPublish = append(hub.onPublish, f)
	return nil
}

func (hub *Events) OnMessage(h MessageHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onMessage = append(hub.onMessage, h)
	return nil
}

func (hub *Events) OnDisconnect(f DisconnectHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onDisconnect = append(hub.onDisconnect, f)
	return nil
}

func (hub *Events) OnInactivate(f InactivateHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onInactivate = append(hub.onInactivate, f)

	return nil
}

func (hub *Events) OnReconnect(f ReconnectHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onReconnect = append(hub.onReconnect, f)
	return nil
}

func (hub *Events) OnChange(f ChangeHandler) {
	hub.mu.Lock()
	hub.onChange = append(hub.onChange, f)
	hub.mu.Unlock()
}

func (hub *Events) Close() error {
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
	hub.onReconnect = nil
	hub.onDisconnect = nil
	hub.onInactivate = nil
	hub.onMessage = nil
	hub.onChange = nil
	hub.mu.Unlock()

	reg.mu.Lock()
	delete(reg.hubs, id)
	reg.mu.Unlock()

	return nil
}

type eventsStorage struct {
	hubs map[int64]*Events
	idSeq int64
	priority Priority
	mu sync.RWMutex
}

func (s *eventsStorage) iterEvents(h func(e *Events)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ev := range s.hubs {
		h(ev)
	}
}

func (s *eventsStorage) new() *Events {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.idSeq++
	id := s.idSeq
	ev := &Events{
		id:           id,
		closed:       false,
		reg:          s,
		mu:           sync.RWMutex{},
	}
	s.hubs[id] = ev
	return ev
}

type sortablePET []*eventsStorage

func (p *sortablePET) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *sortablePET) Push(el interface{}) {
	*p = append(*p, el.(*eventsStorage))
}

func (p *sortablePET) Pop() interface{} {
	el := (*p)[p.Len()-1]
	*p = (*p)[:p.Len()-1]
	return el
}

func (p *sortablePET) Len() int {
	return len(*p)
}

func (p *sortablePET) Less(i, j int) bool {
	return (*p)[i].priority < (*p)[j].priority
}

type eventsRegistry struct {
	tables []*eventsStorage
	mu     sync.RWMutex
}

func (r *eventsRegistry) new(p Priority) *Events {
	r.mu.Lock()
	var sto *eventsStorage
	for _, s := range r.tables {
		if s.priority == p {
			sto = s
			break
		}
	}
	if sto == nil {
		sto = &eventsStorage{
			hubs: map[int64]*Events{},
			idSeq:    0,
			priority: p,
		}
		heap.Push((*sortablePET)(&r.tables), sto)
	}
	r.mu.Unlock()
	return sto.new()
}

func (r *eventsRegistry) emitReconnect(a *App, opts ConnectOptions, c Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Events) {
			e.emitReconnect(wg, a, opts, c)
		})
		wg.Wait()
	}
}

func (r *eventsRegistry) emitMessage(a *App, c Client, data []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Events) {
			e.emitMessage(wg, a, c, data)
		})
		wg.Wait()
	}
}

func (r *eventsRegistry) emitConnect(a *App, opts ConnectOptions, c Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Events) {
			e.emitConnect(wg, a, opts, c)
		})
		wg.Wait()
	}
}

func (r *eventsRegistry) emitDisconnect(a *App, c Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Events) {
			e.emitDisconnect(wg, a, c)
		})
		wg.Wait()
	}
}

func (r *eventsRegistry) emitInactivate(a *App, c Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Events) {
			e.emitInactivate(wg, a, c)
		})
		wg.Wait()
	}
}

func (r *eventsRegistry) emitChange(a *App, c ChangeLog) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Events) {
			e.emitChange(wg, a, c)
		})
		wg.Wait()
	}
}

func newEventsRegistry() *eventsRegistry {
	return &eventsRegistry{}
}