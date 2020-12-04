package pubsub

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	clientid "github.com/RomanIschenko/notify/pubsub/client_id"
	"github.com/google/uuid"
	"sync"
)

type EventHandle struct {
	closer func()
}

func (handle EventHandle) Close() {
	if handle.closer != nil {
		handle.closer()
	}
}

type eventsHub struct {
	ctxHubs map[context.Context]*ContextEvents
	mu sync.RWMutex
}

func newEventsHub() *eventsHub {
	return &eventsHub{
		ctxHubs: map[context.Context]*ContextEvents{},
	}
}

func (hub *eventsHub) ctxHub(ctx context.Context) *ContextEvents {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	h, ok := hub.ctxHubs[ctx]
	if !ok {
		h = newCtxHub()
		hub.ctxHubs[ctx] = h
		if ctx != nil {
			go func(ctx context.Context) {
				select {
				case <-ctx.Done():
					hub.mu.Lock()
					defer hub.mu.Unlock()
					delete(hub.ctxHubs, ctx)
				}
			}(ctx)
		}
	}
	return h
}

func newCtxHub() *ContextEvents {
	return &ContextEvents{
		onChange: map[string]func(arg0 changelog.Log){},
		onPublish: map[string]func(arg0 Publication, arg1 PublishOptions){},
		onSubscribe: map[string]func(arg0 SubscribeOptions, arg1 changelog.Log){},
		onUnsubscribe: map[string]func(arg0 UnsubscribeOptions, arg1 changelog.Log){},
		onConnect: map[string]func(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log){},
		onDisconnect: map[string]func(arg0 DisconnectOptions, arg2 changelog.Log){},
		onInactivate: map[string]func(arg0 clientid.ID){},
	}
}

type ContextEvents struct {
	mu sync.RWMutex
	onChange map[string]func(arg0 changelog.Log)
	onPublish map[string]func(arg0 Publication, arg1 PublishOptions)
	onSubscribe map[string]func(arg0 SubscribeOptions, arg1 changelog.Log)
	onUnsubscribe map[string]func(arg0 UnsubscribeOptions, arg1 changelog.Log)
	onConnect map[string]func(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log)
	onDisconnect map[string]func(arg0 DisconnectOptions, arg2 changelog.Log)
	onInactivate map[string]func(arg0 clientid.ID)
}

func (hub *ContextEvents) emitChange(arg0 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onChange {
		handler(arg0)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitChange(arg0 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitChange(arg0)

	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnChange(f func(arg0 changelog.Log)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onChange[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onChange, uid)
		},
	}
}

func (hub *ContextEvents) emitPublish(arg0 Publication, arg1 PublishOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onPublish {
		handler(arg0, arg1)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitPublish(arg0 Publication, arg1 PublishOptions) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitPublish(arg0, arg1)

	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnPublish(f func(arg0 Publication, arg1 PublishOptions)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onPublish[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onPublish, uid)
		},
	}
}

func (hub *ContextEvents) emitSubscribe(arg0 SubscribeOptions, arg1 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onSubscribe {
		handler(arg0, arg1)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitSubscribe(arg0 SubscribeOptions, arg1 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitSubscribe(arg0, arg1)
		h.emitChange(arg1)
	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnSubscribe(f func(arg0 SubscribeOptions, arg1 changelog.Log)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onSubscribe[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onSubscribe, uid)
		},
	}
}

func (hub *ContextEvents) emitUnsubscribe(arg0 UnsubscribeOptions, arg1 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onUnsubscribe {
		handler(arg0, arg1)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitUnsubscribe(arg0 UnsubscribeOptions, arg1 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitUnsubscribe(arg0, arg1)
		h.emitChange(arg1)
	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnUnsubscribe(f func(arg0 UnsubscribeOptions, arg1 changelog.Log)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onUnsubscribe[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onUnsubscribe, uid)
		},
	}
}

func (hub *ContextEvents) emitConnect(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onConnect {
		handler(arg0, arg1, arg2)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitConnect(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitConnect(arg0, arg1, arg2)
		h.emitChange(arg2)
	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnConnect(f func(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onConnect[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onConnect, uid)
		},
	}
}

func (hub *ContextEvents) emitDisconnect(arg0 DisconnectOptions, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onDisconnect {
		handler(arg0, arg2)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitDisconnect(arg0 DisconnectOptions, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitDisconnect(arg0, arg2)
		h.emitChange(arg2)
	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnDisconnect(f func(arg0 DisconnectOptions, arg2 changelog.Log)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onDisconnect[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onDisconnect, uid)
		},
	}
}

func (hub *ContextEvents) emitInactivate(arg0 clientid.ID) {
	hub.mu.RLock()
	for _, handler := range hub.onInactivate {
		handler(arg0)
	}
	hub.mu.RUnlock()
}
func (hub *eventsHub) emitInactivate(arg0 clientid.ID) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitInactivate(arg0)

	}
	hub.mu.RUnlock()
}
func (hub *ContextEvents) OnInactivate(f func(arg0 clientid.ID)) EventHandle {

	hub.mu.Lock()
	uid := uuid.New().String()
	hub.onInactivate[uid] = f
	hub.mu.Unlock()

	return EventHandle{
		closer: func() {
			hub.mu.Lock()
			defer hub.mu.Unlock()
			delete(hub.onInactivate, uid)
		},
	}
}

