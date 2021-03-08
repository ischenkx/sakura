package notify

import (
	"errors"
	"sync"
	"sync/atomic"
)

var proxyIdSeq = uint64(0)

type Proxy struct {
	connect []func(*App, *ConnectOptions)
	disconnect []func(*App, *DisconnectOptions)
	inactivate []func(*App,  Client)
	subscribe []func(*App, *SubscribeOptions)
	unsubscribe []func(*App, *UnsubscribeOptions)
	publish []func(*App, *PublishOptions)
	reg *proxyRegistry
	closed bool
	id uint64
	mu sync.RWMutex
}

func newProxy() *Proxy {
	return &Proxy{}
}

type proxyRegistry struct {
	mu sync.RWMutex
	hubs map[uint64]*Proxy
}

func (r *proxyRegistry) newHub() *Proxy {
	id := atomic.AddUint64(&proxyIdSeq, 1)
	r.mu.Lock()
	defer r.mu.Unlock()
	h := newProxy()
	r.hubs[id] = h
	return h
}

func (r *proxyRegistry) emitConnect(a *App, opts *ConnectOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.connect {
			f(a, opts)
		}
		h.mu.RUnlock()
	}
}

func (r *proxyRegistry) emitDisconnect(a *App, opts *DisconnectOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.disconnect {
			f(a, opts)
		}
		h.mu.RUnlock()
	}
}

func (r *proxyRegistry) emitInactivate(a *App, c Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.inactivate {
			f(a, c)
		}
		h.mu.RUnlock()
	}
}

func (r *proxyRegistry) emitPublish(a *App, opts *PublishOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.publish {
			f(a, opts)
		}
		h.mu.RUnlock()
	}
}

func (r *proxyRegistry) emitSubscribe(a *App, opts *SubscribeOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.subscribe {
			f(a, opts)
		}
		h.mu.RUnlock()
	}
}

func (r *proxyRegistry) emitUnsubscribe(a *App, opts *UnsubscribeOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.unsubscribe {
			f(a, opts)
		}
		h.mu.RUnlock()
	}
}

func newProxyRegistry() *proxyRegistry {
	return &proxyRegistry{
		mu:   sync.RWMutex{},
		hubs: map[uint64]*Proxy{},
	}
}

func (h *Proxy) OnConnect(f func(*App, *ConnectOptions)) error {
	if f == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return errors.New("proxy: already closed")
	}
	h.connect = append(h.connect, f)
	return nil
}

func (h *Proxy) OnDisconnect(f func(*App, *DisconnectOptions)) error {
	if f == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return errors.New("proxy: already closed")
	}
	h.disconnect = append(h.disconnect, f)
	return nil
}

func (h *Proxy) OnInactivation(f func(*App,  Client)) error {
	if f == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return errors.New("proxy: already closed")
	}
	h.inactivate = append(h.inactivate, f)
	return nil
}

func (h *Proxy) OnSubscribe(f func(*App, *SubscribeOptions)) error {
	if f == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return errors.New("proxy: already closed")
	}
	h.subscribe = append(h.subscribe, f)
	return nil
}

func (h *Proxy) OnUnsubscribe(f func(*App, *UnsubscribeOptions)) error {
	if f == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return errors.New("proxy: already closed")
	}
	h.unsubscribe = append(h.unsubscribe, f)
	return nil
}

func (h *Proxy) OnPublish(f func(*App, *PublishOptions)) error {
	if f == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return errors.New("proxy: already closed")
	}
	h.publish = append(h.publish, f)
	return nil
}

func (h *Proxy) Close() error {
	h.mu.Lock()
	if h.closed {
		h.mu.Unlock()
		return errors.New("proxy: already closed")
	}

	reg, id := h.reg, h.id
	h.reg = nil
	h.closed = true
	h.mu.Unlock()

	reg.mu.Lock()
	delete(reg.hubs, id)
	reg.mu.Unlock()
	return nil
}