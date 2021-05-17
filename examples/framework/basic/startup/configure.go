package startup

import (
	"github.com/ischenkx/notify"
	data2 "github.com/ischenkx/notify/examples/framework/basic/services/data"
	other2 "github.com/ischenkx/notify/examples/framework/basic/services/other"
	"github.com/ischenkx/notify/framework/builder"
	authmock "github.com/ischenkx/notify/pkg/auth/mock"
)

// notify:config
func Configure(builder *builder.Builder) {

	builder.ProvideInjectable("dataService", &data2.Service{
		Data:        "Hello World from data service",
	})

	builder.ProvideInjectable("", &other2.Service{
		Data: "Hello World from other service",
	})

	builder.UseConfig(notify.Config{
		ID:     "my-app",
		Auth:   authmock.New(),
	})
}