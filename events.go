package notify

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub"
	"sync"
)

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
	onChange []func(a *App, log pubsub.ChangeLog)
	onPublish []func(a *App, opts pubsub.PublishOptions)
	onSubscribe []func(a *App, opts pubsub.SubscribeOptions)
	onUnsubscribe []func(a *App, opts pubsub.UnsubscribeOptions)
	onConnect []func(a *App, opts pubsub.ConnectOptions, c pubsub.Client)
	onDisconnect []func(a *App, options pubsub.DisconnectOptions)
	onInactivate []func(a *App, c pubsub.Client)
	onIncomingData []func(a *App, data IncomingData)
}

func newEventsHub() *EventsHub {
	return &EventsHub{}
}

func (hub *EventsHub) emitChange(a *App, log pubsub.ChangeLog) {
	hub.mu.RLock()
	for _, handler := range hub.onChange {
		handler(a, log)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitChange(a *App, log pubsub.ChangeLog) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitChange(a, log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnChange(f func(a *App, l pubsub.ChangeLog)) {
	hub.mu.Lock()
	hub.onChange = append(hub.onChange, f)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitPublish(a *App, opts pubsub.PublishOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onPublish {
		handler(a, opts)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitPublish(a *App, opts pubsub.PublishOptions) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitPublish(a, opts)

	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnPublish(f func(*App, pubsub.PublishOptions)) {

	hub.mu.Lock()
	hub.onPublish = append(hub.onPublish, f)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitIncomingData(a *App, data IncomingData) {
	hub.mu.RLock()
	for _, handler := range hub.onIncomingData {
		handler(a, data)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitIncomingData(a *App, data IncomingData) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitIncomingData(a, data)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnIncomingData(h func(a *App, data IncomingData)) {
	hub.mu.Lock()
	hub.onIncomingData = append(hub.onIncomingData, h)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitSubscribe(a *App, opts pubsub.SubscribeOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onSubscribe {
		handler(a, opts)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitSubscribe(a *App, opts pubsub.SubscribeOptions, log pubsub.ChangeLog) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitSubscribe(a, opts)
		h.emitChange(a, log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnSubscribe(f func(a *App, opts pubsub.SubscribeOptions)) {
	hub.mu.Lock()
	hub.onSubscribe = append(hub.onSubscribe, f)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitUnsubscribe(a *App, opts pubsub.UnsubscribeOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onUnsubscribe {
		handler(a, opts)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitUnsubscribe(a *App, opts pubsub.UnsubscribeOptions, log pubsub.ChangeLog) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitUnsubscribe(a, opts)
		h.emitChange(a, log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnUnsubscribe(f func(a *App, opts pubsub.UnsubscribeOptions)) {

	hub.mu.Lock()
	hub.onUnsubscribe = append(hub.onUnsubscribe, f)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitConnect(a *App, opts pubsub.ConnectOptions, c pubsub.Client) {
	hub.mu.RLock()
	for _, handler := range hub.onConnect {
		handler(a, opts, c)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitConnect(a *App, opts pubsub.ConnectOptions, c pubsub.Client, log pubsub.ChangeLog) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitConnect(a, opts, c)
		h.emitChange(a, log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnConnect(f func(a *App, opts pubsub.ConnectOptions, c pubsub.Client)) {

	hub.mu.Lock()
	hub.onConnect = append(hub.onConnect, f)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitDisconnect(a *App, opts pubsub.DisconnectOptions) {
	hub.mu.RLock()
	for _, handler := range hub.onDisconnect {
		handler(a, opts)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitDisconnect(a *App, opts pubsub.DisconnectOptions, log pubsub.ChangeLog) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitDisconnect(a, opts)
		h.emitChange(a, log)
	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnDisconnect(f func(a *App, opts pubsub.DisconnectOptions)) {

	hub.mu.Lock()
	hub.onDisconnect = append(hub.onDisconnect, f)
	hub.mu.Unlock()
}

func (hub *EventsHub) emitInactivate(a *App, c pubsub.Client) {
	hub.mu.RLock()
	for _, handler := range hub.onInactivate {
		handler(a, c)
	}
	hub.mu.RUnlock()
}
func (hub *eventsRegistry) emitInactivate(a *App, c pubsub.Client) {
	hub.mu.RLock()
	for _, h := range hub.ctxHubs {
		h.emitInactivate(a, c)

	}
	hub.mu.RUnlock()
}
func (hub *EventsHub) OnInactivate(f func(a *App, c pubsub.Client)) {
	hub.mu.Lock()
	hub.onInactivate = append(hub.onInactivate, f)
	hub.mu.Unlock()
}



