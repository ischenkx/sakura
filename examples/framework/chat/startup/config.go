package startup

import (
	chat2 "github.com/ischenkx/notify/examples/framework/chat"
	"github.com/ischenkx/notify/framework/builder"
	"github.com/ischenkx/notify/framework/ioc"
	authmock "github.com/ischenkx/notify/pkg/auth/mock"
)

// notify:configure
func Configure(builder *builder.Builder) {
	builder.AppConfig().Auth = authmock.New()
	builder.IocContainer().AddEntry(ioc.NewEntry(&chat2.Service{Name: "my-service"}, ""))
}

