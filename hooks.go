package notify

import "github.com/RomanIschenko/notify/internal/events"

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

type IncomingEvent struct {
	Name string
	Data interface{}
}

type (
	EmitHandler              func(app *App, event Event)
	ConnectHandler           func(app *App, opts ConnectOptions, client Client)
	ReconnectHandler         func(app *App, opts ConnectOptions, client Client)
	DisconnectHandler        func(app *App, client Client)
	InactivateHandler        func(app *App, client Client)
	EventHandler             func(app *App, client Client, event IncomingEvent)
	ChangeHandler            func(app *App, log ChangeLog)
	BeforeSubscribeHandler   func(*App, *SubscribeOptions)
	BeforeUnsubscribeHandler func(*App, *UnsubscribeOptions)
	BeforeEmitHandler        func(*App, *Event)
	BeforeConnectHandler     func(*App, *ConnectOptions)
	BeforeDisconnectHandler  func(*App, *DisconnectOptions)
)

const (
	receiveHookName           = "Receive"
	emitHookName              = "Emit"
	connectHookName           = "Connect"
	reconnectHookName         = "Reconnect"
	disconnectHookName        = "Disconnect"
	inactivateHookName        = "Inactivate"
	changeHookName            = "Change"
	beforeSubscribeHookName   = "BeforeSubscribe"
	beforeUnsubscribeHookName = "BeforeUnsubscribe"
	beforeEmitHookName        = "BeforePublish"
	beforeConnectHookName     = "BeforeConnect"
	beforeDisconnectHookName  = "BeforeDisconnect"
)

type Hooks struct {
	hub *events.Hub
}

func (h Hooks) OnEmit(f EmitHandler) {
	h.hub.On(emitHookName, f)
}

func (h Hooks) OnEvent(f EventHandler) {
	h.hub.On(receiveHookName, f)
}

func (h Hooks) OnConnect(f ConnectHandler) {
	h.hub.On(connectHookName, f)
}

func (h Hooks) OnDisconnect(f DisconnectHandler) {
	h.hub.On(disconnectHookName, f)
}

func (h Hooks) OnReconnect(f ReconnectHandler) {
	h.hub.On(reconnectHookName, f)
}

func (h Hooks) OnInactivate(f InactivateHandler) {
	h.hub.On(inactivateHookName, f)
}

func (h Hooks) OnChange(f ChangeHandler) {
	h.hub.On(changeHookName, f)
}

func (h Hooks) OnBeforeSubscribe(f BeforeSubscribeHandler) {
	h.hub.On(beforeSubscribeHookName, f)
}

func (h Hooks) OnBeforeUnsubscribe(f BeforeUnsubscribeHandler) {
	h.hub.On(beforeUnsubscribeHookName, f)
}

func (h Hooks) OnBeforeEmit(f BeforeEmitHandler) {
	h.hub.On(beforeEmitHookName, f)
}

func (h Hooks) OnBeforeConnect(f BeforeConnectHandler) {
	h.hub.On(beforeConnectHookName, f)
}

func (h Hooks) OnBeforeDisconnect(f BeforeDisconnectHandler) {
	h.hub.On(beforeDisconnectHookName, f)
}

func (h Hooks) Close() {
	h.hub.Close()
}