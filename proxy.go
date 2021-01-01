package notify

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub"
	"sync"
)

type proxyRegistry struct {
	mu sync.RWMutex
	hubs map[context.Context]*Proxy
}

func (r *proxyRegistry) hub(ctx context.Context) *Proxy {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.hubs[ctx]
	if !ok {
		h = newProxy()
		r.hubs[ctx] = h
		if ctx != nil {
			go func() {
				select {
				case <-ctx.Done():
					r.mu.Lock()
					delete(r.hubs, ctx)
					r.mu.Unlock()
				}
			}()
		}
	}
	return h
}

func (r *proxyRegistry) emitConnect(a *App, opts *pubsub.ConnectOptions) {
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

func (r *proxyRegistry) emitDisconnect(a *App, opts *pubsub.DisconnectOptions) {
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

func (r *proxyRegistry) emitInactivate(a *App, c pubsub.Client) {
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

func (r *proxyRegistry) emitPublish(a *App, opts *pubsub.PublishOptions) {
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

func (r *proxyRegistry) emitSubscribe(a *App, opts *pubsub.SubscribeOptions) {
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

func (r *proxyRegistry) emitUnsubscribe(a *App, opts *pubsub.UnsubscribeOptions) {
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
		hubs: map[context.Context]*Proxy{},
	}
}

type Proxy struct {
	connect []func(*App, *pubsub.ConnectOptions)
	disconnect []func(*App, *pubsub.DisconnectOptions)
	inactivate []func(*App, pubsub.Client)
	subscribe []func(*App, *pubsub.SubscribeOptions)
	unsubscribe []func(*App, *pubsub.UnsubscribeOptions)
	publish []func(*App, *pubsub.PublishOptions)
	mu sync.RWMutex
}

func newProxy() *Proxy {
	return &Proxy{}
}

func (h *Proxy) OnConnect(f func(*App, *pubsub.ConnectOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connect = append(h.connect, f)
}

func (h *Proxy) OnDisconnect(f func(*App, *pubsub.DisconnectOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.disconnect = append(h.disconnect, f)
}

func (h *Proxy) OnInactivation(f func(*App, pubsub.Client)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.inactivate = append(h.inactivate, f)
}

func (h *Proxy) OnSubscribe(f func(*App, *pubsub.SubscribeOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribe = append(h.subscribe, f)
}

func (h *Proxy) OnUnsubscribe(f func(*App, *pubsub.UnsubscribeOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.unsubscribe = append(h.unsubscribe, f)
}

func (h *Proxy) OnPublish(f func(*App, *pubsub.PublishOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.publish = append(h.publish, f)
}