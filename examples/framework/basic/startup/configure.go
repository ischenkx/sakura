package startup

import (
	"github.com/RomanIschenko/notify"
	data2 "github.com/RomanIschenko/notify/examples/framework/basic/services/data"
	other2 "github.com/RomanIschenko/notify/examples/framework/basic/services/other"
	"github.com/RomanIschenko/notify/framework/builder"
	authmock "github.com/RomanIschenko/notify/pkg/auth/mock"
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