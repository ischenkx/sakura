package notify

import "github.com/ischenkx/notify/internal/events"

type IncomingEvent struct {
	Name string
	Data []interface{}
}
type (
	EmitHandler                    = func(*App, Event)
	ConnectHandler                 = func(*App, ConnectOptions, Client)
	ReconnectHandler               = func(*App, ConnectOptions, Client)
	DisconnectHandler              = func(*App, Client)
	InactivateHandler              = func(*App, Client)
	IncomingEventHandler           = func(*App, Client, IncomingEvent)
	ChangeHandler                  = func(*App, ChangeLog)
	ClientSubscribeHandler         = func(*App, Client, SubscribeClientOptions)
	ClientUnsubscribeHandler       = func(*App, Client, UnsubscribeClientOptions)
	UserSubscribeHandler           = func(*App, SubscribeUserOptions)
	UserUnsubscribeHandler         = func(*App, UnsubscribeUserOptions)
	FailedIncomingEventHandler     = func(*App, Client, EventError)
	BeforeClientSubscribeHandler   = func(*App, *SubscribeClientOptions)
	BeforeClientUnsubscribeHandler = func(*App, *UnsubscribeClientOptions)
	BeforeUserSubscribeHandler     = func(*App, *SubscribeUserOptions)
	BeforeUserUnsubscribeHandler   = func(*App, *UnsubscribeUserOptions)
	BeforeEmitHandler              = func(*App, *Event)
	BeforeConnectHandler           = func(*App, *ConnectOptions)
	BeforeDisconnectHandler        = func(*App, *DisconnectOptions)
)


const (
	emitHookName                    = "Emit"
	connectHookName                 = "Connect"
	reconnectHookName               = "Reconnect"
	disconnectHookName              = "Disconnect"
	inactivateHookName              = "Inactivate"
	receiveHookName                 = "Event"
	failedReceiveHookName           = "FailedReceive"
	changeHookName                  = "Change"
	clientSubscribeHookName         = "ClientSubscribe"
	clientUnsubscribeHookName       = "ClientUnsubscribe"
	userSubscribeHookName           = "UserSubscribe"
	userUnsubscribeHookName         = "UserUnsubscribe"
	beforeClientSubscribeHookName   = "BeforeClientSubscribe"
	beforeClientUnsubscribeHookName = "BeforeClientUnsubscribe"
	beforeUserSubscribeHookName     = "BeforeUserSubscribe"
	beforeUserUnsubscribeHookName   = "BeforeUserUnsubscribe"
	beforeEmitHookName              = "BeforeEmit"
	beforeConnectHookName           = "BeforeConnect"
	beforeDisconnectHookName        = "BeforeDisconnect"
)

type Hooks struct {
	hub *events.Hub
}

func (h Hooks) OnEmit(f EmitHandler) {
	h.hub.On(emitHookName, f)
}

func (h Hooks) OnEvent(f IncomingEventHandler) {
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

func (h Hooks) OnFailedEvent(f FailedIncomingEventHandler) {
	h.hub.On(failedReceiveHookName, f)
}

func (h Hooks) OnChange(f ChangeHandler) {
	h.hub.On(changeHookName, f)
}

func (h Hooks) OnUserUnsubscribe(f UserUnsubscribeHandler) {
	h.hub.On(userUnsubscribeHookName, f)
}

func (h Hooks) OnUserSubscribe(f UserSubscribeHandler) {
	h.hub.On(userSubscribeHookName, f)
}

func (h Hooks) OnClientSubscribe(f ClientSubscribeHandler) {
	h.hub.On(clientSubscribeHookName, f)
}

func (h Hooks) OnClientUnsubscribe(f ClientUnsubscribeHandler) {
	h.hub.On(clientUnsubscribeHookName, f)
}

func (h Hooks) OnBeforeClientSubscribe(f BeforeClientSubscribeHandler) {
	h.hub.On(beforeClientSubscribeHookName, f)
}

func (h Hooks) OnBeforeClientUnsubscribe(f BeforeClientUnsubscribeHandler) {
	h.hub.On(beforeClientUnsubscribeHookName, f)
}

func (h Hooks) OnBeforeUserSubscribe(f BeforeUserSubscribeHandler) {
	h.hub.On(beforeUserSubscribeHookName, f)
}

func (h Hooks) OnBeforeUserUnsubscribe(f BeforeUserUnsubscribeHandler) {
	h.hub.On(beforeUserUnsubscribeHookName, f)
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
