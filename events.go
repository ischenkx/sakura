package notify

import (
	"errors"
	"sync"
	"sync/atomic"
)

// TODO
// Events 2.0
// Connect - first connection of the Client
// Reconnect - Client reconnected
// Disconnect - Client disconnected and was deleted
// Inactivate - Client disconnected and is being waited to reconnect
// Message - incoming message
// Change - all changes made in pub/sub system (ChangeLog)
// Publish - new publication
// Subscribe - subscribed
// Unsubscribe - unsubscribed

var idSeq = uint64(0)

type AppEvents struct {
	id uint64
	closed bool
	reg *eventsRegistry
	onPublish []func(a *App, opts PublishOptions)
	onConnect []func(a *App, opts ConnectOptions, c Client)
	onReconnect []func(a *App, opts ConnectOptions, c Client)
	onDisconnect []func(a *App, cs Client)
	onInactivate []func(a *App, c Client)
	onMessage []func(a *App, c Client, message []byte)
	onChange []func(a *App, log ChangeLog)
	mu sync.RWMutex
}

type eventsRegistry struct {
	hubs map[uint64]*AppEvents
	mu   sync.RWMutex
}

func newEventsRegistry() *eventsRegistry {
	return &eventsRegistry{
		hubs: map[uint64]*AppEvents{},
	}
}

func (reg *eventsRegistry) newHub() *AppEvents {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	id := atomic.AddUint64(&idSeq, 1)

	h := &AppEvents{
		id:  id,
		reg: reg,
	}

	reg.hubs[id] = h

	return h
}

func (hub *AppEvents) emitReconnect(a *App, opts ConnectOptions, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onReconnect {
		handler(a, opts, c)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitReconnect(a *App, opts ConnectOptions, c Client) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitReconnect(a, opts, c)
	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnReconnect(f func(a *App, opts ConnectOptions, l Client)) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("AppEvents: closed")
	}

	hub.onReconnect = append(hub.onReconnect, f)
	return nil
}

func (hub *AppEvents) emitPublish(a *App, opts PublishOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onPublish {
		handler(a, opts)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitPublish(a *App, opts PublishOptions) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitPublish(a, opts)

	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnPublish(f func(*App, PublishOptions)) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("AppEvents: closed")
	}

	hub.onPublish = append(hub.onPublish, f)
	return nil
}

func (hub *AppEvents) emitMessage(a *App, c Client, data []byte) {
	hub.mu.RLock()
	for _, handler := range hub.onMessage {
		handler(a, c, data)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitMessage(a *App, c Client, data []byte) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitMessage(a, c, data)
	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnMessage(h func(a *App, c Client, data []byte)) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("AppEvents: closed")
	}

	hub.onMessage = append(hub.onMessage, h)
	return nil
}

func (hub *AppEvents) emitConnect(a *App, opts ConnectOptions, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onConnect {
		handler(a, opts, c)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitConnect(a *App, opts ConnectOptions, c Client) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitConnect(a, opts, c)
	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnConnect(f func(a *App, opts ConnectOptions, c Client)) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("AppEvents: closed")
	}
	hub.onConnect = append(hub.onConnect, f)
	return nil
}

func (hub *AppEvents) emitDisconnect(a *App, client Client) {
	hub.mu.RLock()
	for _, handler := range hub.onDisconnect {
		handler(a, client)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitDisconnect(a *App, client Client) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitDisconnect(a, client)
	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnDisconnect(f func(a *App, clients Client)) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return errors.New("AppEvents: closed")
	}
	hub.onDisconnect = append(hub.onDisconnect, f)
	return nil
}

func (hub *AppEvents) emitInactivate(a *App, c Client) {
	hub.mu.RLock()
	for _, handler := range hub.onInactivate {
		handler(a, c)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitInactivate(a *App, c Client) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitInactivate(a, c)

	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnInactivate(f func(a *App, c Client)) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if hub.closed {
		return errors.New("AppEvents: closed")
	}

	hub.onInactivate = append(hub.onInactivate, f)

	return nil
}

func (hub *AppEvents) emitChange(a *App, c ChangeLog) {
	hub.mu.RLock()
	for _, handler := range hub.onChange {
		handler(a, c)
	}
	hub.mu.RUnlock()
}

func (reg *eventsRegistry) emitChange(a *App, c ChangeLog) {
	reg.mu.RLock()
	for _, h := range reg.hubs {
		h.emitChange(a, c)

	}
	reg.mu.RUnlock()
}

func (hub *AppEvents) OnChange(f func(a *App, c ChangeLog)) {
	hub.mu.Lock()
	hub.onChange = append(hub.onChange, f)
	hub.mu.Unlock()
}

func (hub *AppEvents) Close() error {
	hub.mu.Lock()
	if hub.closed {
		hub.mu.Unlock()
		return errors.New("AppEvents: already closed")
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


