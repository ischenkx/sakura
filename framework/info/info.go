package info

import (
	"context"
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/framework/ioc"
)

const ImportPath = "github.com/ischenkx/notify/framework/info"

type Info struct {
	ctx context.Context
	app *notify.App
	ioc *ioc.Container
}

func (info *Info) App() *notify.App {
	return info.app
}

func (info *Info) IocContainer() *ioc.Container {
	return info.ioc
}

func (info *Info) Context() context.Context {
	return info.ctx
}

func New(ctx context.Context, ioc *ioc.Container, app *notify.App) Info {
	return Info{
		ctx: ctx,
		app: app,
		ioc: ioc,
	}
}
