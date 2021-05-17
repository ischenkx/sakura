package runtime

import (
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

type Hook struct {
	Event string
	Handler interface{}
	Priority notify.Priority
}

type HookContainer struct {
	app *notify.App
	hooks map[notify.Priority]notify.Hooks
}

func (c HookContainer) getOrCreate(p notify.Priority) notify.Hooks {
	h, ok := c.hooks[p]

	if !ok {
		h = c.app.Hooks(p)
		c.hooks[p] = h
	}

	return h
}

func (c HookContainer) Close() {
	for _, hooks := range c.hooks {
		hooks.Close()
	}
}

func (c HookContainer) Init(rawHooks []Hook) {
	for _, h := range rawHooks {
		hooks := c.getOrCreate(h.Priority)
		event, handler := h.Event, h.Handler
		switch event {
		case "emit":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.EmitHandler)).Elem()) {
				panic("failed to match emit hook")
				continue
			}
			hooks.OnEmit(handler.(notify.EmitHandler))
		case "connect":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.ConnectHandler)).Elem()) {
				panic("failed to match connect hook")
				continue
			}
			hooks.OnConnect(handler.(notify.ConnectHandler))
		case "reconnect":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.ReconnectHandler)).Elem()) {
				panic("failed to match reconnect hook")
				continue
			}
			hooks.OnReconnect(handler.(notify.ReconnectHandler))
		case "disconnect":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.DisconnectHandler)).Elem()) {
				panic("failed to match disconnect hook")
				continue
			}
			hooks.OnDisconnect(handler.(notify.DisconnectHandler))
		case "inactivate":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.InactivateHandler)).Elem()) {
				panic("failed to match inactivate hook")
				continue
			}
			hooks.OnInactivate(handler.(notify.InactivateHandler))
		case "incomingEvent":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.IncomingEventHandler)).Elem()) {
				panic("failed to match event hook")
				continue
			}
			hooks.OnEvent(handler.(notify.IncomingEventHandler))
		case "failedEvent":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.FailedIncomingEventHandler)).Elem()) {
				panic("failed to match failedReceive hook")
				continue
			}
			hooks.OnFailedEvent(handler.(notify.FailedIncomingEventHandler))
		case "change":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.ChangeHandler)).Elem()) {
				panic("failed to match change hook")
				continue
			}
			hooks.OnChange(handler.(notify.ChangeHandler))
		case "clientSubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.ClientSubscribeHandler)).Elem()) {
				panic("failed to match clientSubscribe hook")
				continue
			}
			hooks.OnClientSubscribe(handler.(notify.ClientSubscribeHandler))
		case "clientUnsubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.ClientUnsubscribeHandler)).Elem()) {
				panic("failed to match clientUnsubscribe hook")
				continue
			}
			hooks.OnClientUnsubscribe(handler.(notify.ClientUnsubscribeHandler))
		case "userSubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.UserSubscribeHandler)).Elem()) {
				panic("failed to match userSubscribe hook")
				continue
			}
			hooks.OnUserSubscribe(handler.(notify.UserSubscribeHandler))
		case "userUnsubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.UserUnsubscribeHandler)).Elem()) {
				panic("failed to match userUnsubscribe hook")
				continue
			}
			hooks.OnUserUnsubscribe(handler.(notify.UserUnsubscribeHandler))
		case "beforeClientSubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeClientSubscribeHandler)).Elem()) {
				panic("failed to match beforeClientSubscribe hook")
				continue
			}
			hooks.OnBeforeClientSubscribe(handler.(notify.BeforeClientSubscribeHandler))
		case "beforeClientUnsubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeClientUnsubscribeHandler)).Elem()) {
				panic("failed to match beforeClientUnsubscribe hook")
				continue
			}
			hooks.OnBeforeClientUnsubscribe(handler.(notify.BeforeClientUnsubscribeHandler))
		case "beforeUserSubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeUserSubscribeHandler)).Elem()) {
				panic("failed to match beforeUserSubscribe hook")
				continue
			}
			hooks.OnBeforeUserSubscribe(handler.(notify.BeforeUserSubscribeHandler))
		case "beforeUserUnsubscribe":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeUserUnsubscribeHandler)).Elem()) {
				panic("failed to match beforeUserUnsubscribe hook")
				continue
			}
			hooks.OnBeforeUserUnsubscribe(handler.(notify.BeforeUserUnsubscribeHandler))
		case "beforeEmit":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeEmitHandler)).Elem()) {
				panic("failed to match beforeEmit hook")
				continue
			}
			hooks.OnBeforeEmit(handler.(notify.BeforeEmitHandler))
		case "beforeConnect":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeConnectHandler)).Elem()) {
				panic("failed to match beforeConnect hook")
				continue
			}
			hooks.OnBeforeConnect(handler.(notify.BeforeConnectHandler))
		case "beforeDisconnect":
			if !matchHandlerAndHandlerType(handler, reflect.TypeOf(new(notify.BeforeDisconnectHandler)).Elem()) {
				panic("failed to match beforeDisconnect hook")
				continue
			}
			hooks.OnBeforeDisconnect(handler.(notify.BeforeDisconnectHandler))
		}
	}
}

func matchHandlerAndHandlerType(handler interface{}, t reflect.Type) bool {

	ht := reflect.TypeOf(handler)

	if ht.Kind() != reflect.Func || t.Kind() != reflect.Func {
		return false
	}

	if ht.NumIn() != t.NumIn() || ht.NumOut() != t.NumOut(){
		return false
	}

	for i := 0; i < ht.NumOut(); i++ {
		if !utils.CompareTypes(ht.Out(i), t.Out(i)) {
			return false
		}
	}

	for i := 0; i < ht.NumIn(); i++ {
		if !utils.CompareTypes(ht.In(i), t.In(i)) {
			return false
		}
	}

	return true
}


func NewHookContainer(app *notify.App) HookContainer {
	container := HookContainer{
		app:   app,
		hooks: map[notify.Priority]notify.Hooks{},
	}
	return container
}
