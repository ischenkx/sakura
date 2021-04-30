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
	PublishHandler    func(app *App, opts PublishOptions)
	ConnectHandler    func(app *App, opts ConnectOptions, client Client)
	ReconnectHandler  func(app *App, opts ConnectOptions, client Client)
	DisconnectHandler func(app *App, client Client)
	InactivateHandler func(app *App, client Client)
	MessageHandler    func(app *App, client Client, message []byte)
	ChangeHandler     func(app *App, log ChangeLog)
	BeforeSubscribeHandler func(*App, *SubscribeOptions)
	BeforeUnsubscribeHandler func(*App, *UnsubscribeOptions)
	BeforePublishHandler func(*App, *PublishOptions)
	BeforeConnectHandler func(*App, *ConnectOptions)
	BeforeDisconnectHandler func(*App, *DisconnectOptions)
)

type Hooks struct {
	id           int64
	closed       bool
	reg          *eventsStorage
	onPublish    []PublishHandler
	onConnect    []ConnectHandler
	onReconnect  []ReconnectHandler
	onDisconnect []DisconnectHandler
	onInactivate []InactivateHandler
	onChange     []ChangeHandler
	onBeforeUnsubscribe  []UnsubscribeProxyHandler
	onBeforeSubscribe  []SubscribeProxyHandler
	onBeforePublish    []PublishProxyHandler
	onBeforeConnect    []ConnectProxyHandler
	onBeforeDisconnect []DisconnectProxyHandler
	mu           sync.RWMutex
}

func (hub *Hooks) emitReconnect(wg *sync.WaitGroup, a *App, opts ConnectOptions, c Client) {
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

func (hub *Hooks) emitPublish(wg *sync.WaitGroup, a *App, opts PublishOptions) {
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

func (hub *Hooks) emitConnect(wg *sync.WaitGroup, a *App, opts ConnectOptions, c Client) {
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

func (hub *Hooks) emitDisconnect(wg *sync.WaitGroup, a *App, c Client) {
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

func (hub *Hooks) emitInactivate(wg *sync.WaitGroup, a *App, c Client) {
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

func (hub *Hooks) emitChange(wg *sync.WaitGroup, a *App, c ChangeLog) {
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

func (hub *Hooks) OnConnect(f ConnectHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onConnect = append(hub.onConnect, f)
	return nil
}

func (hub *Hooks) OnPublish(f PublishHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onPublish = append(hub.onPublish, f)
	return nil
}

func (hub *Hooks) OnDisconnect(f DisconnectHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onDisconnect = append(hub.onDisconnect, f)
	return nil
}

func (hub *Hooks) OnInactivate(f InactivateHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onInactivate = append(hub.onInactivate, f)

	return nil
}

func (hub *Hooks) OnReconnect(f ReconnectHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onReconnect = append(hub.onReconnect, f)
	return nil
}

func (hub *Hooks) OnChange(f ChangeHandler) {
	hub.mu.Lock()
	hub.onChange = append(hub.onChange, f)
	hub.mu.Unlock()
}

func (hub *Hooks) emitBeforeSubscribe(wg *sync.WaitGroup, a *App, opts *SubscribeOptions) {
	hub.mu.RLock()
	for _, handler := range hub {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Hooks) emitBeforeUnsubscribe(wg *sync.WaitGroup, a *App, opts *UnsubscribeOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onBeforeUnsubscribe {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Proxy) emitBeforePublish(wg *sync.WaitGroup, a *App, opts *PublishOptions) {
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

func (hub *Proxy) emitBeforeConnect(wg *sync.WaitGroup, a *App, opts *ConnectOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onConnect {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Proxy) emitBeforeDisconnect(wg *sync.WaitGroup, a *App, opts *DisconnectOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onDisconnect {
		h := handler
		wg.Add(1)
		go func() {
			h(a, opts)
			wg.Done()
		}()
	}
	hub.mu.RUnlock()
}

func (hub *Hooks) OnBeforeSubscribe(f SubscribeProxyHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onSubscribe = append(hub.onSubscribe, f)
	return nil
}

func (hub *Hooks) OnBeforeConnect(f ConnectProxyHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onConnect = append(hub.onConnect, f)
	return nil
}

func (hub *Hooks) OnBeforePublish(f PublishProxyHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}

	hub.onPublish = append(hub.onPublish, f)
	return nil
}

func (hub *Hooks) OnBeforeDisconnect(f DisconnectProxyHandler) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("Proxy: closed")
	}
	hub.onDisconnect = append(hub.onDisconnect, f)
	return nil
}

func (hub *Hooks) Close() error {
	hub.mu.Lock()
	if hub.closed {
		hub.mu.Unlock()
		return errors.New("proxy: already closed")
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
	hub.onChange = nil
	hub.mu.Unlock()

	reg.mu.Lock()
	delete(reg.hubs, id)
	reg.mu.Unlock()

	return nil
}

type eventsStorage struct {
	hubs map[int64]*Hooks
	idSeq int64
	priority Priority
	mu sync.RWMutex
}

func (s *eventsStorage) iterEvents(h func(e *Hooks)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ev := range s.hubs {
		h(ev)
	}
}

func (s *eventsStorage) new() *Hooks {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.idSeq++
	id := s.idSeq
	ev := &Hooks{
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

func (r *eventsRegistry) new(p Priority) *Hooks {
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
			hubs: map[int64]*Hooks{},
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
		t.iterEvents(func(e *Hooks) {
			e.emitReconnect(wg, a, opts, c)
		})
		wg.Wait()
	}
}

func (r *eventsRegistry) emitMessage(a *App, c Client, data []byte) {

}

func (r *eventsRegistry) emitConnect(a *App, opts ConnectOptions, c Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tables {
		wg := &sync.WaitGroup{}
		t.iterEvents(func(e *Hooks) {
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
		t.iterEvents(func(e *Hooks) {
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
		t.iterEvents(func(e *Hooks) {
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
		t.iterEvents(func(e *Hooks) {
			e.emitChange(wg, a, c)
		})
		wg.Wait()
	}
}

func newEventsRegistry() *eventsRegistry {
	return &eventsRegistry{}
}