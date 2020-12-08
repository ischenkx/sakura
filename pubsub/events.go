package pubsub

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub/changelog"
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

type eventsRegistry struct {
	ctxHubs map[context.Context]*EventsHub
	mu sync.RWMutex
}

func newEventsRegistry() *eventsRegistry {
	return &eventsRegistry{
		ctxHubs: map[context.Context]*EventsHub{},
	}
}

func (hub *eventsRegistry) hub(ctx context.Context) *EventsHub {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	h, ok := hub.ctxHubs[ctx]
	if !ok {
		h = newEventsHub()
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

type EventsHub struct {
	mu sync.RWMutex
	onChange map[string]func(changelog.Log)
	onPublish map[string]func(PublishOptions)
	onSubscribe map[string]func(SubscribeOptions, changelog.Log)
	onUnsubscribe map[string]func(UnsubscribeOptions, changelog.Log)
	onConnect map[string]func(ConnectOptions, *Client, changelog.Log)
	onDisconnect map[string]func(DisconnectOptions, changelog.Log)
	onInactivate map[string]func(string)
}

func newEventsHub() *EventsHub {
	return &EventsHub{
		onChange: map[string]func(changelog.Log){},
		onPublish: map[string]func(PublishOptions){},
		onSubscribe: map[string]func(SubscribeOptions, changelog.Log){},
		onUnsubscribe: map[string]func(UnsubscribeOptions, changelog.Log){},
		onConnect: map[string]func(ConnectOptions, *Client, changelog.Log){},
		onDisconnect: map[string]func(DisconnectOptions, changelog.Log){},
		onInactivate: map[string]func(string){},
	}
}

func (hub *EventsHub) emitChange(log changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onChange {
		handler(log)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitChange(log changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitChange(log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnChange(f func(changelog.Log)) EventHandle {

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

func (hub *EventsHub) emitPublish(opts PublishOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onPublish {
		handler(opts)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitPublish(opts PublishOptions) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitPublish(opts)

	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnPublish(f func(PublishOptions)) EventHandle {

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

func (hub *EventsHub) emitSubscribe(opts SubscribeOptions, log changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onSubscribe {
		handler(opts, log)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitSubscribe(opts SubscribeOptions, log changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitSubscribe(opts, log)
		h.emitChange(log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnSubscribe(f func(arg0 SubscribeOptions, arg1 changelog.Log)) EventHandle {

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

func (hub *EventsHub) emitUnsubscribe(arg0 UnsubscribeOptions, arg1 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onUnsubscribe {
		handler(arg0, arg1)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitUnsubscribe(arg0 UnsubscribeOptions, arg1 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitUnsubscribe(arg0, arg1)
		h.emitChange(arg1)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnUnsubscribe(f func(arg0 UnsubscribeOptions, arg1 changelog.Log)) EventHandle {

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

func (hub *EventsHub) emitConnect(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onConnect {
		handler(arg0, arg1, arg2)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitConnect(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitConnect(arg0, arg1, arg2)
		h.emitChange(arg2)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnConnect(f func(arg0 ConnectOptions, arg1 *Client, arg2 changelog.Log)) EventHandle {

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

func (hub *EventsHub) emitDisconnect(arg0 DisconnectOptions, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, handler := range hub.onDisconnect {
		handler(arg0, arg2)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitDisconnect(arg0 DisconnectOptions, arg2 changelog.Log) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitDisconnect(arg0, arg2)
		h.emitChange(arg2)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnDisconnect(f func(arg0 DisconnectOptions, arg2 changelog.Log)) EventHandle {

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

func (hub *EventsHub) emitInactivate(arg0 string) {
	hub.mu.RLock()
	for _, handler := range hub.onInactivate {
		handler(arg0)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitInactivate(arg0 string) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitInactivate(arg0)

	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnInactivate(f func(arg0 string)) EventHandle {

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

