package simplehandler

import (
	"github.com/ischenkx/notify/examples/framework/app/services"
	"github.com/ischenkx/notify/framework/builder"
	authmock "github.com/ischenkx/notify/pkg/auth/mock"
)

// notify:configure
func Configurator(b *builder.Builder) {
	b.AppConfig().ID = "simple-app"
	b.AppConfig().Auth = authmock.New()

	b.IocContainer().Inject(&services.UniqueIDProvider{})
}
