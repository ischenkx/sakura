package main

import (
	notify "github.com/ischenkx/notify"
	gateway "github.com/ischenkx/notify/examples/framework/app/gateway"
	hooks "github.com/ischenkx/notify/examples/framework/app/hooks"
	startup "github.com/ischenkx/notify/examples/framework/app/startup"
	builder "github.com/ischenkx/notify/framework/builder"
	hmapper "github.com/ischenkx/notify/framework/hmapper"
	info "github.com/ischenkx/notify/framework/info"
	ioc "github.com/ischenkx/notify/framework/ioc"
	runtime "github.com/ischenkx/notify/framework/runtime"
	reflect "reflect"
)

var configurators = []interface{}{
	startup.Configurator,
}

var hookList = []runtime.Hook{
	{Event: "connect", Handler: hooks.ConnectHook, Priority: 65536},
	{Event: "failedIncomingEvent", Handler: hooks.FailedEventHook, Priority: 65536},
}

var consumers = []ioc.Consumer{}

type handler struct {
	example interface{}

	prefix string

	eventsMapping map[string]string
}

var handlers = []handler{
	{example: &gateway.Handler{}, eventsMapping: map[string]string{"joinStrings": "OnJoinStrings", "message": "OnMessage", "request": "OnRequest"}, prefix: ""},
}

func main() {
	iocContainer := ioc.New()
	appBuilder := builder.New(iocContainer)
	runtime.Configure(appBuilder, configurators...)
	app := notify.New(*appBuilder.AppConfig())
	iocContainer.Inject(app)
	handlersMapper := hmapper.New(app)
	hookContainer := runtime.NewHookContainer(app)
	hookContainer.Init(hookList)
	defer hookContainer.Close()
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
	runtimeInfo := info.New(appBuilder.Context(), iocContainer, app)
	runtime.Start(runtimeInfo, startup.Start)
}
