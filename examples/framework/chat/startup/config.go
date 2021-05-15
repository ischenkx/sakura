package startup

import (
	"github.com/RomanIschenko/notify"
	chat2 "github.com/RomanIschenko/notify/examples/framework/chat"
	"github.com/RomanIschenko/notify/framework/builder"
	authmock "github.com/RomanIschenko/notify/pkg/auth/mock"
)

// notify:config
func Configure(builder *builder.Builder) {
	builder.SetAppConfig(notify.Config{
		ID:   "my-app",
		Auth: authmock.New(),
	})
	builder.Inject(&chat2.Service{Name: "my-service"})
}

