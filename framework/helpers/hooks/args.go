package hooks

import (
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

type EmitArgs struct {
	App   *notify.App
	Event notify.Event
}
type ConnectArgs struct {
	App     *notify.App
	Options notify.ConnectOptions
	Client  notify.Client
}
type ReconnectArgs struct {
	App     *notify.App
	Options notify.ConnectOptions
	Client  notify.Client
}
type DisconnectArgs struct {
	App    *notify.App
	Client notify.Client
}
type InactivateArgs struct {
	App    *notify.App
	Client notify.Client
}
type IncomingEventArgs struct {
	App    *notify.App
	Client notify.Client
	Event  notify.IncomingEvent
}
type ChangeArgs struct {
	App *notify.App
	Log notify.ChangeLog
}
type ClientSubscribeArgs struct {
	App     *notify.App
	Client  notify.Client
	Options notify.SubscribeClientOptions
}
type ClientUnsubscribeArgs struct {
	App     *notify.App
	Client  notify.Client
	Options notify.UnsubscribeClientOptions
}
type UserSubscribeArgs struct {
	App     *notify.App
	Options notify.SubscribeUserOptions
}
type UserUnsubscribeArgs struct {
	App     *notify.App
	Options notify.UnsubscribeUserOptions
}
type FailedIncomingEventArgs struct {
	App        *notify.App
	Client     notify.Client
	EventError notify.EventError
}
type BeforeClientSubscribeArgs struct {
	App     *notify.App
	Options *notify.SubscribeClientOptions
}
type BeforeClientUnsubscribeArgs struct {
	App     *notify.App
	Options *notify.UnsubscribeClientOptions
}
type BeforeUserSubscribeArgs struct {
	App     *notify.App
	Options *notify.SubscribeUserOptions
}
type BeforeUserUnsubscribeArgs struct {
	App     *notify.App
	Options *notify.UnsubscribeUserOptions
}
type BeforeEmitArgs struct {
	App   *notify.App
	Event *notify.Event
}
type BeforeConnectArgs struct {
	App     *notify.App
	Options *notify.ConnectOptions
}
type BeforeDisconnectArgs struct {
	App     *notify.App
	Options *notify.DisconnectOptions
}

func TryToTransformHandler(h interface{}) (interface{}, bool) {
	val := reflect.ValueOf(h)

	htyp := val.Type()

	if htyp.Kind() != reflect.Func {
		return nil, false
	}

	if htyp.NumIn() != 1 {
		return nil, false
	}

	firstParamType := htyp.In(0)

	switch {
	case utils.CompareTypes(firstParamType, reflect.TypeOf(EmitArgs{})):
		return func(arg0 *notify.App, arg1 notify.Event, ) {
			args := EmitArgs{App: arg0, Event: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(ConnectArgs{})):
		return func(arg0 *notify.App, arg1 notify.ConnectOptions, arg2 notify.Client, ) {
			args := ConnectArgs{App: arg0, Options: arg1, Client: arg2}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(ReconnectArgs{})):
		return func(arg0 *notify.App, arg1 notify.ConnectOptions, arg2 notify.Client, ) {
			args := ReconnectArgs{App: arg0, Options: arg1, Client: arg2}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(DisconnectArgs{})):
		return func(arg0 *notify.App, arg1 notify.Client, ) {
			args := DisconnectArgs{App: arg0, Client: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(InactivateArgs{})):
		return func(arg0 *notify.App, arg1 notify.Client, ) {
			args := InactivateArgs{App: arg0, Client: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(IncomingEventArgs{})):
		return func(arg0 *notify.App, arg1 notify.Client, arg2 notify.IncomingEvent, ) {
			args := IncomingEventArgs{App: arg0, Client: arg1, Event: arg2}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(ChangeArgs{})):
		return func(arg0 *notify.App, arg1 notify.ChangeLog, ) {
			args := ChangeArgs{App: arg0, Log: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(ClientSubscribeArgs{})):
		return func(arg0 *notify.App, arg1 notify.Client, arg2 notify.SubscribeClientOptions, ) {
			args := ClientSubscribeArgs{App: arg0, Client: arg1, Options: arg2}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(ClientUnsubscribeArgs{})):
		return func(arg0 *notify.App, arg1 notify.Client, arg2 notify.UnsubscribeClientOptions, ) {
			args := ClientUnsubscribeArgs{App: arg0, Client: arg1, Options: arg2}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(UserSubscribeArgs{})):
		return func(arg0 *notify.App, arg1 notify.SubscribeUserOptions, ) {
			args := UserSubscribeArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(UserUnsubscribeArgs{})):
		return func(arg0 *notify.App, arg1 notify.UnsubscribeUserOptions, ) {
			args := UserUnsubscribeArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(FailedIncomingEventArgs{})):
		return func(arg0 *notify.App, arg1 notify.Client, arg2 notify.EventError) {
			args := FailedIncomingEventArgs{App: arg0, Client: arg1, EventError: arg2}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeClientSubscribeArgs{})):
		return func(arg0 *notify.App, arg1 *notify.SubscribeClientOptions, ) {
			args := BeforeClientSubscribeArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeClientUnsubscribeArgs{})):
		return func(arg0 *notify.App, arg1 *notify.UnsubscribeClientOptions, ) {
			args := BeforeClientUnsubscribeArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeUserSubscribeArgs{})):
		return func(arg0 *notify.App, arg1 *notify.SubscribeUserOptions, ) {
			args := BeforeUserSubscribeArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeUserUnsubscribeArgs{})):
		return func(arg0 *notify.App, arg1 *notify.UnsubscribeUserOptions, ) {
			args := BeforeUserUnsubscribeArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeEmitArgs{})):
		return func(arg0 *notify.App, arg1 *notify.Event, ) {
			args := BeforeEmitArgs{App: arg0, Event: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeConnectArgs{})):
		return func(arg0 *notify.App, arg1 *notify.ConnectOptions, ) {
			args := BeforeConnectArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	case utils.CompareTypes(firstParamType, reflect.TypeOf(BeforeDisconnectArgs{})):
		return func(arg0 *notify.App, arg1 *notify.DisconnectOptions, ) {
			args := BeforeDisconnectArgs{App: arg0, Options: arg1}
			val.Call([]reflect.Value{reflect.ValueOf(args)})
		}, true

	}
	return nil, false
}
