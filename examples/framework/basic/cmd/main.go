package main

import (
	xhxKQFDaFpLS "github.com/ischenkx/notify"
	"github.com/ischenkx/notify/examples/framework/basic"
	"github.com/ischenkx/notify/examples/framework/basic/services/data"
	"github.com/ischenkx/notify/examples/framework/basic/services/other"
	"github.com/ischenkx/notify/examples/framework/basic/startup"
	jFbcXoEFfRsW "github.com/ischenkx/notify/framework/builder"
	AjWwhTHctcuA "github.com/ischenkx/notify/framework/hmapper"
	XVlBzgbaiCMR "github.com/ischenkx/notify/framework/ioc"
	"github.com/ischenkx/notify/framework/runtime"
	wekrBEmfdzdc "reflect"
)

func main() {
	ioc := XVlBzgbaiCMR.New()
	builder := jFbcXoEFfRsW.New(ioc)
	configurators := []interface{}{startup.Configure}
	runtime.Configure(builder, configurators...)
	app := xhxKQFDaFpLS.New(*builder.Config())
	ioc.Add(XVlBzgbaiCMR.NewEntry(app, ""))
	hmapper := AjWwhTHctcuA.New(app)
	entities := []entity{
		{
			value:   &basic.Handler{},
			mapping: map[string]string{"app": "", "s1": "", "s2": ""},
		},

		{
			value:   &data.Service{},
			mapping: map[string]string{"OtherService": ""},
		},

		{
			value:   &other.Service{},
			mapping: map[string]string{"DataService": "dataService"},
		},
	}
	for _, ent := range entities {
		ioc.InitConsumer(ent.value, ent.mapping)
	}
	handler0, ok := ioc.FindConsumer(wekrBEmfdzdc.TypeOf(&basic.Handler{}))
	if !ok {
		handler0 = &basic.Handler{}
	}
	handler0_mapping := map[string]string{}
	handler0_mapping["message"] = "HandleMessage"
	hmapper.AddHandler(handler0, "basic", handler0_mapping)
	starter := startup.Start
	runtime.Start(app, starter)
}

type entity struct {
	mapping map[string]string

	value interface{}
}
