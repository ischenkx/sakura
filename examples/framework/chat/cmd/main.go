package main

import (
	notify "github.com/ischenkx/notify"
	chat "github.com/ischenkx/notify/examples/framework/chat"
	chat1 "github.com/ischenkx/notify/examples/framework/chat/chat"
	hooks "github.com/ischenkx/notify/examples/framework/chat/hooks"
	startup "github.com/ischenkx/notify/examples/framework/chat/startup"
	builder "github.com/ischenkx/notify/framework/builder"
	hmapper "github.com/ischenkx/notify/framework/hmapper"
	info "github.com/ischenkx/notify/framework/info"
	ioc "github.com/ischenkx/notify/framework/ioc"
	runtime "github.com/ischenkx/notify/framework/runtime"
	reflect "reflect"
)

var configurators = []interface{}{
	startup.Configure,
}

var hookList = []runtime.Hook{
	{Event: "beforeConnect", Handler: hooks.OnBeforeConnect, Priority: 0},
	{Event: "reconnect", Handler: hooks.ReconnectHook, Priority: 65536},
	{Event: "connect", Handler: hooks.onConnect, Priority: 123},
}

var consumers = []ioc.Consumer{
	{Object: &chat.Chat{}, DependencyMapping: map[string]string{"app": "", "service": ""}},
	{Object: &chat.Handler0{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler1{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler10{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler11{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler12{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler13{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler14{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler15{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler16{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler17{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler18{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler19{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler2{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler20{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler21{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler22{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler23{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler24{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler25{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler26{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler27{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler28{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler29{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler3{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler30{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler4{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler5{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler6{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler7{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler8{}, DependencyMapping: map[string]string{"app": ""}},
	{Object: &chat.Handler9{}, DependencyMapping: map[string]string{"app": ""}},
}

type handler struct {
	prefix string

	eventsMapping map[string]string

	example interface{}
}

var handlers = []handler{
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

var initializers = []interface{}{
	startup.Initializer,
}

func main() {
	iocContainer := ioc.New()
	appBuilder := builder.New(iocContainer)
	runtime.Configure(appBuilder, configurators...)
	app := notify.New(*appBuilder.AppConfig())
	iocContainer.AddEntry(ioc.NewEntry(app, ""))
	handlersMapper := hmapper.New(app)
	hookContainer := runtime.NewHookContainer(app)
	hookContainer.Init(hookList)
	defer hookContainer.Close()
	runtimeInfo := info.New(appBuilder.Context(), iocContainer, app)
	for _, c := range consumers {
		iocContainer.AddConsumer(c)
	}
	for _, h := range handlers {
		handlerValue := h.example
		if existingHandler, ok := iocContainer.FindConsumer(reflect.TypeOf(h.example)); ok {
			handlerValue = existingHandler
		}
		handlersMapper.AddHandler(handlerValue, h.prefix, h.eventsMapping)
	}
	runtime.Initialize(runtimeInfo, initializers...)
	runtime.Start(app, startup.Start)
}
