package swirl

import (
	"container/heap"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Priority int

// -----------------------
// A LOT OF GENERATED CODE
// -----------------------

type localAppEvents struct {
	priority           Priority
	registry           *eventsRegistry
	id                 string
	mu                 sync.RWMutex
	_Emit              []func(*App, EmitOptions)
	_Event             []func(*App, Client, EventOptions)
	_Connect           []func(*App, ConnectOptions, Client)
	_Disconnect        []func(*App, Client)
	_Reconnect         []func(*App, ConnectOptions, Client)
	_Inactivate        []func(*App, Client)
	_Error             []func(*App, error)
	_Change            []func(*App, ChangeLog)
	_ClientSubscribe   []func(*App, Client, string, int64)
	_UserSubscribe     []func(*App, User, string, int64)
	_ClientUnsubscribe []func(*App, Client, string, int64)
	_UserUnsubscribe   []func(*App, User, string, int64)
}

func (evs *localAppEvents) callEmit(arg0 EmitOptions) {
	evs.mu.RLock()
	handlers := make([]func(*App, EmitOptions), len(evs._Emit))
	copy(handlers, evs._Emit)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0)
	}
}
func (evs *localAppEvents) OnEmit(f func(*App, EmitOptions)) {
	evs.mu.Lock()
	evs._Emit = append(evs._Emit, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callEvent(arg0 Client, arg1 EventOptions) {
	evs.mu.RLock()
	handlers := make([]func(*App, Client, EventOptions), len(evs._Event))
	copy(handlers, evs._Event)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1)
	}
}
func (evs *localAppEvents) OnEvent(f func(*App, Client, EventOptions)) {
	evs.mu.Lock()
	evs._Event = append(evs._Event, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callConnect(arg0 ConnectOptions, arg1 Client) {
	evs.mu.RLock()
	handlers := make([]func(*App, ConnectOptions, Client), len(evs._Connect))
	copy(handlers, evs._Connect)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1)
	}
}
func (evs *localAppEvents) OnConnect(f func(*App, ConnectOptions, Client)) {
	evs.mu.Lock()
	evs._Connect = append(evs._Connect, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callDisconnect(arg0 Client) {
	evs.mu.RLock()
	handlers := make([]func(*App, Client), len(evs._Disconnect))
	copy(handlers, evs._Disconnect)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0)
	}
}
func (evs *localAppEvents) OnDisconnect(f func(*App, Client)) {
	evs.mu.Lock()
	evs._Disconnect = append(evs._Disconnect, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callReconnect(arg0 ConnectOptions, arg1 Client) {
	evs.mu.RLock()
	handlers := make([]func(*App, ConnectOptions, Client), len(evs._Reconnect))
	copy(handlers, evs._Reconnect)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1)
	}
}
func (evs *localAppEvents) OnReconnect(f func(*App, ConnectOptions, Client)) {
	evs.mu.Lock()
	evs._Reconnect = append(evs._Reconnect, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callInactivate(arg0 Client) {
	evs.mu.RLock()
	handlers := make([]func(*App, Client), len(evs._Inactivate))
	copy(handlers, evs._Inactivate)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0)
	}
}
func (evs *localAppEvents) OnInactivate(f func(*App, Client)) {
	evs.mu.Lock()
	evs._Inactivate = append(evs._Inactivate, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callError(arg0 error) {
	evs.mu.RLock()
	handlers := make([]func(*App, error), len(evs._Error))
	copy(handlers, evs._Error)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0)
	}
}
func (evs *localAppEvents) OnError(f func(*App, error)) {
	evs.mu.Lock()
	evs._Error = append(evs._Error, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callChange(arg0 ChangeLog) {
	evs.mu.RLock()
	handlers := make([]func(*App, ChangeLog), len(evs._Change))
	copy(handlers, evs._Change)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0)
	}
}
func (evs *localAppEvents) OnChange(f func(*App, ChangeLog)) {
	evs.mu.Lock()
	evs._Change = append(evs._Change, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callClientSubscribe(arg0 Client, arg1 string, arg2 int64) {
	evs.mu.RLock()
	handlers := make([]func(*App, Client, string, int64), len(evs._ClientSubscribe))
	copy(handlers, evs._ClientSubscribe)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1, arg2)
	}
}
func (evs *localAppEvents) OnClientSubscribe(f func(*App, Client, string, int64)) {
	evs.mu.Lock()
	evs._ClientSubscribe = append(evs._ClientSubscribe, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callUserSubscribe(arg0 User, arg1 string, arg2 int64) {
	evs.mu.RLock()
	handlers := make([]func(*App, User, string, int64), len(evs._UserSubscribe))
	copy(handlers, evs._UserSubscribe)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1, arg2)
	}
}
func (evs *localAppEvents) OnUserSubscribe(f func(*App, User, string, int64)) {
	evs.mu.Lock()
	evs._UserSubscribe = append(evs._UserSubscribe, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callClientUnsubscribe(arg0 Client, arg1 string, arg2 int64) {
	evs.mu.RLock()
	handlers := make([]func(*App, Client, string, int64), len(evs._ClientUnsubscribe))
	copy(handlers, evs._ClientUnsubscribe)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1, arg2)
	}
}
func (evs *localAppEvents) OnClientUnsubscribe(f func(*App, Client, string, int64)) {
	evs.mu.Lock()
	evs._ClientUnsubscribe = append(evs._ClientUnsubscribe, f)
	evs.mu.Unlock()
}
func (evs *localAppEvents) callUserUnsubscribe(arg0 User, arg1 string, arg2 int64) {
	evs.mu.RLock()
	handlers := make([]func(*App, User, string, int64), len(evs._UserUnsubscribe))
	copy(handlers, evs._UserUnsubscribe)
	evs.mu.RUnlock()
	for _, h := range handlers {
		h(evs.registry.app, arg0, arg1, arg2)
	}
}
func (evs *localAppEvents) OnUserUnsubscribe(f func(*App, User, string, int64)) {
	evs.mu.Lock()
	evs._UserUnsubscribe = append(evs._UserUnsubscribe, f)
	evs.mu.Unlock()
}

type eventsRegistry struct {
	app              *App
	mu               sync.RWMutex
	appEventHandlers []*localAppEvents
	localEntities    *localEntitiesEventsRegistry
}

func (reg *eventsRegistry) callEmit(arg0 EmitOptions) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callEmit(arg0)
	}
}
func (reg *eventsRegistry) callEvent(arg0 Client, arg1 EventOptions) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callEvent(arg0, arg1)
	}
}
func (reg *eventsRegistry) callConnect(arg0 ConnectOptions, arg1 Client) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callConnect(arg0, arg1)
	}
}
func (reg *eventsRegistry) callDisconnect(arg0 Client) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callDisconnect(arg0)
	}
}
func (reg *eventsRegistry) callReconnect(arg0 ConnectOptions, arg1 Client) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callReconnect(arg0, arg1)
	}
}
func (reg *eventsRegistry) callInactivate(arg0 Client) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callInactivate(arg0)
	}
}
func (reg *eventsRegistry) callError(arg0 error) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callError(arg0)
	}
}
func (reg *eventsRegistry) callChange(arg0 ChangeLog) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callChange(arg0)
	}
}
func (reg *eventsRegistry) callClientSubscribe(arg0 Client, arg1 string, arg2 int64) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callClientSubscribe(arg0, arg1, arg2)
	}
}
func (reg *eventsRegistry) callUserSubscribe(arg0 User, arg1 string, arg2 int64) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callUserSubscribe(arg0, arg1, arg2)
	}
}
func (reg *eventsRegistry) callClientUnsubscribe(arg0 Client, arg1 string, arg2 int64) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callClientUnsubscribe(arg0, arg1, arg2)
	}
}
func (reg *eventsRegistry) callUserUnsubscribe(arg0 User, arg1 string, arg2 int64) {
	reg.mu.RLock()
	handlers := make([]*localAppEvents, len(reg.appEventHandlers))
	copy(handlers, reg.appEventHandlers)
	reg.mu.RUnlock()
	for _, h := range handlers {
		h.callUserUnsubscribe(arg0, arg1, arg2)
	}
}

type localTopicsEventsRegistry struct {
	mu                 sync.RWMutex
	_presence          map[string]string
	_topicHandlers     map[string][]string
	_ClientSubscribe   map[string][]func(Client)
	_UserSubscribe     map[string][]func(User)
	_ClientUnsubscribe map[string][]func(Client)
	_UserUnsubscribe   map[string][]func(User)
	_Emit              map[string][]func(EventOptions)
}

func (reg *localTopicsEventsRegistry) handleClientSubscribe(uid string, f func(Client)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._ClientSubscribe[uid] = append(reg._ClientSubscribe[uid], f)
}
func (reg *localTopicsEventsRegistry) callClientSubscribe(topicID string, arg0 Client) {
	reg.mu.RLock()
	uids, ok := reg._topicHandlers[topicID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(Client)
	for _, uid := range uids {
		if hs, ok := reg._ClientSubscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localTopicsEventsRegistry) handleUserSubscribe(uid string, f func(User)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._UserSubscribe[uid] = append(reg._UserSubscribe[uid], f)
}
func (reg *localTopicsEventsRegistry) callUserSubscribe(topicID string, arg0 User) {
	reg.mu.RLock()
	uids, ok := reg._topicHandlers[topicID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(User)
	for _, uid := range uids {
		if hs, ok := reg._UserSubscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localTopicsEventsRegistry) handleClientUnsubscribe(uid string, f func(Client)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._ClientUnsubscribe[uid] = append(reg._ClientUnsubscribe[uid], f)
}
func (reg *localTopicsEventsRegistry) callClientUnsubscribe(topicID string, arg0 Client) {
	reg.mu.RLock()
	uids, ok := reg._topicHandlers[topicID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(Client)
	for _, uid := range uids {
		if hs, ok := reg._ClientUnsubscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localTopicsEventsRegistry) handleUserUnsubscribe(uid string, f func(User)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._UserUnsubscribe[uid] = append(reg._UserUnsubscribe[uid], f)
}
func (reg *localTopicsEventsRegistry) callUserUnsubscribe(topicID string, arg0 User) {
	reg.mu.RLock()
	uids, ok := reg._topicHandlers[topicID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(User)
	for _, uid := range uids {
		if hs, ok := reg._UserUnsubscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localTopicsEventsRegistry) handleEmit(uid string, f func(EventOptions)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Emit[uid] = append(reg._Emit[uid], f)
}
func (reg *localTopicsEventsRegistry) callEmit(topicID string, arg0 EventOptions) {
	reg.mu.RLock()
	uids, ok := reg._topicHandlers[topicID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(EventOptions)
	for _, uid := range uids {
		if hs, ok := reg._Emit[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localTopicsEventsRegistry) createEvents(topicID string) localTopicEvents {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	uid := uuid.New().String()

	reg._presence[uid] = topicID
	reg._topicHandlers[topicID] = append(reg._topicHandlers[topicID], uid)

	return localTopicEvents{reg: reg, topicID: topicID, id: uid}
}
func (reg *localTopicsEventsRegistry) deleteEvents(uid string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	topicID, ok := reg._presence[uid]
	if !ok {
		return
	}
	handlers := reg._topicHandlers[topicID]

	for i := 0; i < len(handlers); i++ {
		if handlers[i] == uid {
			handlers[i] = handlers[len(handlers)-1]
			handlers = handlers[:len(handlers)-1]
			break
		}
	}
	reg._topicHandlers[topicID] = handlers
	delete(reg._presence, uid)
	delete(reg._ClientSubscribe, uid)
	delete(reg._UserSubscribe, uid)
	delete(reg._ClientUnsubscribe, uid)
	delete(reg._UserUnsubscribe, uid)
	delete(reg._Emit, uid)
}

type localClientsEventsRegistry struct {
	mu              sync.RWMutex
	_presence       map[string]string
	_clientHandlers map[string][]string
	_Emit           map[string][]func(EventOptions)
	_Event          map[string][]func(EventOptions)
	_Unsubscribe    map[string][]func(string, int64)
	_Subscribe      map[string][]func(string, int64)
	_Reconnect      map[string][]func(int64)
	_Disconnect     map[string][]func(int64)
}

func (reg *localClientsEventsRegistry) handleEmit(uid string, f func(EventOptions)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Emit[uid] = append(reg._Emit[uid], f)
}
func (reg *localClientsEventsRegistry) callEmit(clientID string, arg0 EventOptions) {
	reg.mu.RLock()
	uids, ok := reg._clientHandlers[clientID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(EventOptions)
	for _, uid := range uids {
		if hs, ok := reg._Emit[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localClientsEventsRegistry) handleEvent(uid string, f func(EventOptions)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Event[uid] = append(reg._Event[uid], f)
}
func (reg *localClientsEventsRegistry) callEvent(clientID string, arg0 EventOptions) {
	reg.mu.RLock()
	uids, ok := reg._clientHandlers[clientID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(EventOptions)
	for _, uid := range uids {
		if hs, ok := reg._Event[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localClientsEventsRegistry) handleUnsubscribe(uid string, f func(string, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Unsubscribe[uid] = append(reg._Unsubscribe[uid], f)
}
func (reg *localClientsEventsRegistry) callUnsubscribe(clientID string, arg0 string, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._clientHandlers[clientID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(string, int64)
	for _, uid := range uids {
		if hs, ok := reg._Unsubscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localClientsEventsRegistry) handleSubscribe(uid string, f func(string, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Subscribe[uid] = append(reg._Subscribe[uid], f)
}
func (reg *localClientsEventsRegistry) callSubscribe(clientID string, arg0 string, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._clientHandlers[clientID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(string, int64)
	for _, uid := range uids {
		if hs, ok := reg._Subscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localClientsEventsRegistry) handleReconnect(uid string, f func(int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Reconnect[uid] = append(reg._Reconnect[uid], f)
}
func (reg *localClientsEventsRegistry) callReconnect(clientID string, arg0 int64) {
	reg.mu.RLock()
	uids, ok := reg._clientHandlers[clientID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(int64)
	for _, uid := range uids {
		if hs, ok := reg._Reconnect[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localClientsEventsRegistry) handleDisconnect(uid string, f func(int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Disconnect[uid] = append(reg._Disconnect[uid], f)
}

func (reg *localClientsEventsRegistry) callDisconnect(clientID string, arg0 int64) {
	reg.mu.RLock()
	uids, ok := reg._clientHandlers[clientID]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(int64)
	for _, uid := range uids {
		if hs, ok := reg._Disconnect[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localClientsEventsRegistry) createEvents(clientID string) localClientEvents {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	uid := uuid.New().String()

	reg._presence[uid] = clientID
	reg._clientHandlers[clientID] = append(reg._clientHandlers[clientID], uid)

	return localClientEvents{reg: reg, clientID: clientID, id: uid}
}
func (reg *localClientsEventsRegistry) deleteEvents(uid string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	clientID, ok := reg._presence[uid]
	if !ok {
		return
	}
	handlers := reg._clientHandlers[clientID]

	for i := 0; i < len(handlers); i++ {
		if handlers[i] == uid {
			handlers[i] = handlers[len(handlers)-1]
			handlers = handlers[:len(handlers)-1]
			break
		}
	}
	reg._clientHandlers[clientID] = handlers
	delete(reg._presence, uid)
	delete(reg._Emit, uid)
	delete(reg._Event, uid)
	delete(reg._Unsubscribe, uid)
	delete(reg._Subscribe, uid)
	delete(reg._Reconnect, uid)
	delete(reg._Disconnect, uid)
}

type localUsersEventsRegistry struct {
	mu                sync.RWMutex
	_presence         map[string]string
	_handlers         map[string][]string
	_Emit             map[string][]func(EventOptions)
	_Event            map[string][]func(EventOptions)
	_Unsubscribe      map[string][]func(string, int64)
	_Subscribe        map[string][]func(string, int64)
	_ClientConnect    map[string][]func(Client, int64)
	_ClientReconnect  map[string][]func(Client, int64)
	_ClientDisconnect map[string][]func(string, int64)
}

func (reg *localUsersEventsRegistry) handleEmit(uid string, f func(EventOptions)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Emit[uid] = append(reg._Emit[uid], f)
}
func (reg *localUsersEventsRegistry) callEmit(__id string, arg0 EventOptions) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(EventOptions)
	for _, uid := range uids {
		if hs, ok := reg._Emit[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localUsersEventsRegistry) handleEvent(uid string, f func(EventOptions)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Event[uid] = append(reg._Event[uid], f)
}
func (reg *localUsersEventsRegistry) callEvent(__id string, arg0 EventOptions) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(EventOptions)
	for _, uid := range uids {
		if hs, ok := reg._Event[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0)
	}
}
func (reg *localUsersEventsRegistry) handleUnsubscribe(uid string, f func(string, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Unsubscribe[uid] = append(reg._Unsubscribe[uid], f)
}
func (reg *localUsersEventsRegistry) callUnsubscribe(__id string, arg0 string, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(string, int64)
	for _, uid := range uids {
		if hs, ok := reg._Unsubscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localUsersEventsRegistry) handleSubscribe(uid string, f func(string, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._Subscribe[uid] = append(reg._Subscribe[uid], f)
}
func (reg *localUsersEventsRegistry) callSubscribe(__id string, arg0 string, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(string, int64)
	for _, uid := range uids {
		if hs, ok := reg._Subscribe[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localUsersEventsRegistry) handleClientConnect(uid string, f func(Client, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._ClientConnect[uid] = append(reg._ClientConnect[uid], f)
}
func (reg *localUsersEventsRegistry) callClientConnect(__id string, arg0 Client, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(Client, int64)
	for _, uid := range uids {
		if hs, ok := reg._ClientConnect[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localUsersEventsRegistry) handleClientReconnect(uid string, f func(Client, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._ClientReconnect[uid] = append(reg._ClientReconnect[uid], f)
}
func (reg *localUsersEventsRegistry) callClientReconnect(__id string, arg0 Client, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(Client, int64)
	for _, uid := range uids {
		if hs, ok := reg._ClientReconnect[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localUsersEventsRegistry) handleClientDisconnect(uid string, f func(string, int64)) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	_, ok := reg._presence[uid]
	if !ok {
		return
	}
	reg._ClientDisconnect[uid] = append(reg._ClientDisconnect[uid], f)
}
func (reg *localUsersEventsRegistry) callClientDisconnect(__id string, arg0 string, arg1 int64) {
	reg.mu.RLock()
	uids, ok := reg._handlers[__id]
	if !ok {
		reg.mu.RUnlock()
		return
	}
	var handlers []func(string, int64)
	for _, uid := range uids {
		if hs, ok := reg._ClientDisconnect[uid]; ok {
			handlers = append(handlers, hs...)
		}
	}
	reg.mu.RUnlock()

	for _, h := range handlers {
		h(arg0, arg1)
	}
}
func (reg *localUsersEventsRegistry) createEvents(__id string) localUserEvents {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	uid := uuid.New().String()

	reg._presence[uid] = __id
	reg._handlers[__id] = append(reg._handlers[__id], uid)

	return localUserEvents{reg: reg, userID: __id, id: uid}
}
func (reg *localUsersEventsRegistry) deleteEvents(uid string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	__id, ok := reg._presence[uid]
	if !ok {
		return
	}
	handlers := reg._handlers[__id]

	for i := 0; i < len(handlers); i++ {
		if handlers[i] == uid {
			handlers[i] = handlers[len(handlers)-1]
			handlers = handlers[:len(handlers)-1]
			break
		}
	}
	reg._handlers[__id] = handlers
	delete(reg._presence, uid)
	delete(reg._Emit, uid)
	delete(reg._Event, uid)
	delete(reg._Unsubscribe, uid)
	delete(reg._Subscribe, uid)
	delete(reg._ClientConnect, uid)
	delete(reg._ClientReconnect, uid)
	delete(reg._ClientDisconnect, uid)
}

// -----------------
// HAND-WRITTEN PART
// -----------------

type localEntitiesEventsRegistry struct {
	users   *localUsersEventsRegistry
	clients *localClientsEventsRegistry
	topics  *localTopicsEventsRegistry
}

type localTopicEvents struct {
	reg         *localTopicsEventsRegistry
	id, topicID string
}

func (l localTopicEvents) OnClientSubscribe(f func(Client)) {
	l.reg.handleClientSubscribe(l.id, f)

}

func (l localTopicEvents) OnUserSubscribe(f func(User)) {
	l.reg.handleUserSubscribe(l.id, f)

}

func (l localTopicEvents) OnClientUnsubscribe(f func(Client)) {
	l.reg.handleClientUnsubscribe(l.id, f)

}

func (l localTopicEvents) OnUserUnsubscribe(f func(User)) {
	l.reg.handleUserUnsubscribe(l.id, f)

}

func (l localTopicEvents) OnEmit(f func(options EventOptions)) {
	l.reg.handleEmit(l.id, f)

}

func (l localTopicEvents) Close() {
	l.reg.deleteEvents(l.id)

}

type localUserEvents struct {
	reg        *localUsersEventsRegistry
	id, userID string
}

func (l localUserEvents) OnEmit(f func(EventOptions)) {
	l.reg.handleEmit(l.id, f)

}

func (l localUserEvents) OnEvent(f func(EventOptions)) {
	l.reg.handleEvent(l.id, f)

}

func (l localUserEvents) OnUnsubscribe(f func(topic string, ts int64)) {
	l.reg.handleUnsubscribe(l.id, f)

}

func (l localUserEvents) OnSubscribe(f func(topic string, ts int64)) {
	l.reg.handleSubscribe(l.id, f)

}

func (l localUserEvents) OnClientConnect(f func(Client, int64)) {
	l.reg.handleClientConnect(l.id, f)

}

func (l localUserEvents) OnClientReconnect(f func(Client, int64)) {
	l.reg.handleClientReconnect(l.id, f)

}

func (l localUserEvents) OnClientDisconnect(f func(id string, ts int64)) {
	l.reg.handleClientDisconnect(l.id, f)
}

func (l localUserEvents) Close() {
	l.reg.deleteEvents(l.id)
}

type localClientEvents struct {
	reg          *localClientsEventsRegistry
	id, clientID string
}

func (l localClientEvents) OnEmit(f func(EventOptions)) {
	l.reg.handleEmit(l.id, f)
}

func (l localClientEvents) OnEvent(f func(EventOptions)) {
	l.reg.handleEvent(l.id, f)
}

func (l localClientEvents) OnUnsubscribe(f func(topic string, ts int64)) {
	l.reg.handleUnsubscribe(l.id, f)
}

func (l localClientEvents) OnSubscribe(f func(topic string, ts int64)) {
	l.reg.handleSubscribe(l.id, f)
}

func (l localClientEvents) OnReconnect(f func(int64)) {
	l.reg.handleReconnect(l.id, f)
}

func (l localClientEvents) OnDisconnect(f func(ts int64)) {
	l.reg.handleDisconnect(l.id, f)
}

func (l localClientEvents) Close() {
	l.reg.deleteEvents(l.id)
}

func (evs *localAppEvents) Close() {
	evs.registry.deleteHub(evs.id)
}

type prioritySortedHubList []*localAppEvents

func (p *prioritySortedHubList) Len() int {
	return len(*p)
}

func (p *prioritySortedHubList) Less(i, j int) bool {
	return (*p)[i].priority < (*p)[j].priority
}

func (p *prioritySortedHubList) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *prioritySortedHubList) Push(x interface{}) {
	*p = append(*p, x.(*localAppEvents))
}

func (p *prioritySortedHubList) Pop() interface{} {
	last := (*p)[p.Len()-1]
	*p = (*p)[:p.Len()-1]
	return last
}

func (reg *eventsRegistry) newHub(priority Priority) *localAppEvents {
	appEvs := &localAppEvents{
		priority: priority,
		registry: reg,
		id:       uuid.New().String(),
	}

	reg.mu.Lock()
	heap.Push((*prioritySortedHubList)(&reg.appEventHandlers), appEvs)
	reg.mu.Unlock()

	return appEvs
}

func (reg *eventsRegistry) deleteHub(id string) {
	reg.mu.Lock()
	for i, el := range reg.appEventHandlers {
		if el.id == id {
			heap.Remove((*prioritySortedHubList)(&reg.appEventHandlers), i)
			break
		}
	}
	reg.mu.Unlock()
}

func (reg *eventsRegistry) forApp(priority Priority) AppEvents {
	return reg.newHub(priority)
}

func (reg *eventsRegistry) forClient(id string) ClientEvents {
	return reg.localEntities.clients.createEvents(id)
}

func (reg *eventsRegistry) forUser(id string) UserEvents {
	return reg.localEntities.users.createEvents(id)
}

func (reg *eventsRegistry) forTopic(id string) TopicEvents {
	return reg.localEntities.topics.createEvents(id)
}

func (reg *localClientsEventsRegistry) deleteByClientID(id string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if uids, ok := reg._clientHandlers[id]; ok {
		delete(reg._clientHandlers, id)
		for _, uid := range uids {
			delete(reg._presence, uid)
			delete(reg._Emit, uid)
			delete(reg._Event, uid)
			delete(reg._Unsubscribe, uid)
			delete(reg._Subscribe, uid)
			delete(reg._Reconnect, uid)
			delete(reg._Disconnect, uid)
		}
	}
}

func (reg *eventsRegistry) initBindings() {
	hub := reg.newHub(-1)

	hub.OnClientSubscribe(func(app *App, client Client, s string, i int64) {
		reg.localEntities.clients.callSubscribe(client.ID(), s, i)
		reg.localEntities.topics.callClientSubscribe(s, client)
	})

	hub.OnClientUnsubscribe(func(app *App, client Client, s string, i int64) {
		reg.localEntities.clients.callUnsubscribe(client.ID(), s, i)
		reg.localEntities.topics.callClientUnsubscribe(s, client)
	})

	hub.OnUserSubscribe(func(app *App, user User, s string, i int64) {
		reg.localEntities.users.callSubscribe(user.ID(), s, i)
		reg.localEntities.topics.callUserSubscribe(s, user)

	})

	hub.OnUserUnsubscribe(func(app *App, user User, s string, i int64) {
		reg.localEntities.users.callUnsubscribe(user.ID(), s, i)
		reg.localEntities.topics.callUserUnsubscribe(s, user)
	})

	hub.OnEmit(func(app *App, options EmitOptions) {
		for _, id := range options.Users {
			reg.localEntities.users.callEmit(id, options.EventOptions)
		}

		for _, id := range options.Clients {
			reg.localEntities.clients.callEmit(id, options.EventOptions)

		}

		for _, id := range options.Topics {
			reg.localEntities.topics.callEmit(id, options.EventOptions)

		}
	})

	hub.OnEvent(func(app *App, client Client, options EventOptions) {
		reg.localEntities.clients.callEvent(client.ID(), options)
		if client.User().ID() != "" {
			reg.localEntities.users.callEvent(client.User().ID(), options)
		}
	})

	hub.OnConnect(func(app *App, options ConnectOptions, client Client) {
		if client.User().ID() != "" {
			reg.localEntities.users.callClientConnect(client.User().ID(), client, options.TimeStamp)
		}
	})

	hub.OnReconnect(func(app *App, options ConnectOptions, client Client) {
		reg.localEntities.clients.callReconnect(client.ID(), options.TimeStamp)
		if client.User().ID() != "" {
			reg.localEntities.users.callClientReconnect(client.User().ID(), client, options.TimeStamp)
		}
	})

	hub.OnDisconnect(func(app *App, client Client) {
		reg.localEntities.clients.callDisconnect(client.ID(), time.Now().UnixNano())
		if client.User().ID() != "" {
			reg.localEntities.users.callClientDisconnect(client.User().ID(), client.ID(), time.Now().UnixNano())
		}
	})
}

func newEventsRegistry(app *App) *eventsRegistry {
	reg := &eventsRegistry{
		app:              app,
		mu:               sync.RWMutex{},
		appEventHandlers: []*localAppEvents{},
		localEntities: &localEntitiesEventsRegistry{
			users: &localUsersEventsRegistry{
				mu:                sync.RWMutex{},
				_presence:         map[string]string{},
				_handlers:         map[string][]string{},
				_Emit:             map[string][]func(EventOptions){},
				_Event:            map[string][]func(EventOptions){},
				_Unsubscribe:      map[string][]func(string, int64){},
				_Subscribe:        map[string][]func(string, int64){},
				_ClientConnect:    map[string][]func(Client, int64){},
				_ClientReconnect:  map[string][]func(Client, int64){},
				_ClientDisconnect: map[string][]func(string, int64){},
			},
			clients: &localClientsEventsRegistry{
				mu:              sync.RWMutex{},
				_presence:       map[string]string{},
				_clientHandlers: map[string][]string{},
				_Emit:           map[string][]func(EventOptions){},
				_Event:          map[string][]func(EventOptions){},
				_Unsubscribe:    map[string][]func(string, int64){},
				_Subscribe:      map[string][]func(string, int64){},
				_Reconnect:      map[string][]func(int64){},
				_Disconnect:     map[string][]func(int64){},
			},
			topics: &localTopicsEventsRegistry{
				mu:                 sync.RWMutex{},
				_presence:          map[string]string{},
				_topicHandlers:     map[string][]string{},
				_ClientSubscribe:   map[string][]func(Client){},
				_UserSubscribe:     map[string][]func(User){},
				_ClientUnsubscribe: map[string][]func(Client){},
				_UserUnsubscribe:   map[string][]func(User){},
				_Emit:              map[string][]func(EventOptions){},
			},
		},
	}
	reg.initBindings()
	return reg
}
