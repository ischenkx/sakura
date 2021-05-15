package main

import (
	notify "github.com/RomanIschenko/notify"
	chat "github.com/RomanIschenko/notify/examples/framework/chat"
	chat1 "github.com/RomanIschenko/notify/examples/framework/chat/chat"
	hooks "github.com/RomanIschenko/notify/examples/framework/chat/hooks"
	startup "github.com/RomanIschenko/notify/examples/framework/chat/startup"
	builder "github.com/RomanIschenko/notify/framework/builder"
	hmapper "github.com/RomanIschenko/notify/framework/hmapper"
	ioc "github.com/RomanIschenko/notify/framework/ioc"
	runtime "github.com/RomanIschenko/notify/framework/runtime"
	reflect "reflect"
)

func main() {
	iocContainer := ioc.New()
	appBuilder := builder.New(iocContainer)
	configurators := []interface{}{startup.Configure}
	runtime.Configure(appBuilder, configurators...)
	app := notify.New(*appBuilder.AppConfig())
	iocContainer.Add(ioc.NewEntry(app, ""))
	handlersMapper := hmapper.New(app)
	hookContainer := runtime.NewHookContainer(app)
	hookContainer.Init([]runtime.Hook{{Event: "beforeConnect", Handler: hooks.OnBeforeConnect}, {Event: "reconnect", Handler: hooks.ReconnectHook}})
	defer hookContainer.Close()
	consumers := []consumer{
		{value: &chat.Chat{}, mapping: map[string]string{"app": "", "service": ""}},
		{value: &chat.Handler0{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler1{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler10{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler11{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler12{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler13{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler14{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler15{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler16{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler17{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler18{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler19{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler2{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler20{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler21{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler22{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler23{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler24{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler25{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler26{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler27{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler28{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler29{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler3{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler30{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler4{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler5{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler6{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler7{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler8{}, mapping: map[string]string{"app": ""}},
		{value: &chat.Handler9{}, mapping: map[string]string{"app": ""}},
	}
	for _, c := range consumers {
		iocContainer.InitConsumer(c.value, c.mapping)
	}
	handlers := []handler{
		{example: &chat.Chat{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: ""},
		{example: &chat.Handler0{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler0"},
		{example: &chat.Handler1{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler1"},
		{example: &chat.Handler10{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler10"},
		{example: &chat.Handler11{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler11"},
		{example: &chat.Handler12{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler12"},
		{example: &chat.Handler13{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler13"},
		{example: &chat.Handler14{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler14"},
		{example: &chat.Handler15{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler15"},
		{example: &chat.Handler16{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler16"},
		{example: &chat.Handler17{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler17"},
		{example: &chat.Handler18{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler18"},
		{example: &chat.Handler19{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler19"},
		{example: &chat.Handler2{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler2"},
		{example: &chat.Handler20{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler20"},
		{example: &chat.Handler21{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler21"},
		{example: &chat.Handler22{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler22"},
		{example: &chat.Handler23{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler23"},
		{example: &chat.Handler24{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler24"},
		{example: &chat.Handler25{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler25"},
		{example: &chat.Handler26{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler26"},
		{example: &chat.Handler27{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler27"},
		{example: &chat.Handler28{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler28"},
		{example: &chat.Handler29{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler29"},
		{example: &chat.Handler3{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler3"},
		{example: &chat.Handler30{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler30"},
		{example: &chat.Handler4{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler4"},
		{example: &chat.Handler5{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler5"},
		{example: &chat.Handler6{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler6"},
		{example: &chat.Handler7{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler7"},
		{example: &chat.Handler8{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler8"},
		{example: &chat.Handler9{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "handler9"},
		{example: &chat1.Handler{}, eventsMapping: map[string]string{"message": "HandleMessage"}, prefix: "chat"},
	}
	for _, h := range handlers {
		existingH, ok := iocContainer.FindConsumer(reflect.TypeOf(h.example))
		var handlerValue interface{}
		if ok {
			handlerValue = existingH
		} else {
			handlerValue = h.example
		}
		handlersMapper.AddHandler(handlerValue, h.prefix, h.eventsMapping)
	}
	runtime.Start(app, startup.Start)
}

type consumer struct {
	value interface{}

	mapping map[string]string
}

type handler struct {
	eventsMapping map[string]string

	example interface{}

	prefix string
}
