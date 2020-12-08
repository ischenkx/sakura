package pubsub

import (
	"context"
	"sync"
)

type middlewareRegistry struct {
	mu sync.RWMutex
	hubs map[context.Context]*MiddlewareHub
}

func (r *middlewareRegistry) hub(ctx context.Context) *MiddlewareHub {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.hubs[ctx]
	if !ok {
		h = newMiddlewareHub()
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

func (r *middlewareRegistry) emitConnect(opts *ConnectOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.connect {
			f(opts)
		}
		h.mu.RUnlock()
	}
}

func (r *middlewareRegistry) emitDisconnect(opts *DisconnectOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.disconnect {
			f(opts)
		}
		h.mu.RUnlock()
	}
}

func (r *middlewareRegistry) emitInactivate(c *Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.inactivate {
			f(c)
		}
		h.mu.RUnlock()
	}
}

func (r *middlewareRegistry) emitPublish(opts *PublishOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.publish {
			f(opts)
		}
		h.mu.RUnlock()
	}
}

func (r *middlewareRegistry) emitSubscribe(opts *SubscribeOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.subscribe {
			f(opts)
		}
		h.mu.RUnlock()
	}
}

func (r *middlewareRegistry) emitUnsubscribe(opts *UnsubscribeOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hubs {
		h.mu.RLock()
		for _, f := range h.unsubscribe {
			f(opts)
		}
		h.mu.RUnlock()
	}
}

func newMiddlewareRegistry() *middlewareRegistry {
	return &middlewareRegistry{
		mu:   sync.RWMutex{},
		hubs: map[context.Context]*MiddlewareHub{},
	}
}

type MiddlewareHub struct {
	connect []func(*ConnectOptions)
	disconnect []func(*DisconnectOptions)
	inactivate []func(*Client)
	subscribe []func(*SubscribeOptions)
	unsubscribe []func(*UnsubscribeOptions)
	publish []func(*PublishOptions)
	mu sync.RWMutex
}

func newMiddlewareHub() *MiddlewareHub {
	return &MiddlewareHub{}
}

func (h *MiddlewareHub) OnConnect(f func(*ConnectOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connect = append(h.connect, f)
}

func (h *MiddlewareHub) OnDisconnect(f func(*DisconnectOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.disconnect = append(h.disconnect, f)
}

func (h *MiddlewareHub) OnInactivation(f func(*Client)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.inactivate = append(h.inactivate, f)
}

func (h *MiddlewareHub) OnSubscribe(f func(*SubscribeOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribe = append(h.subscribe, f)
}

func (h *MiddlewareHub) OnUnsubscribe(f func(*UnsubscribeOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.unsubscribe = append(h.unsubscribe, f)
}

func (h *MiddlewareHub) OnPublish(f func(*PublishOptions)) {
	if f == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.publish = append(h.publish, f)
}