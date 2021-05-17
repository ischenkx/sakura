package builder

import (
	"context"
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/framework/ioc"
	"github.com/google/uuid"
)

const ImportPath = "github.com/ischenkx/notify/framework/builder"

type LabelOption struct {
	Name string
}

type Builder struct {
	appConfig notify.Config
	ioc *ioc.Container
	ctx context.Context
}

func (b *Builder) AppConfig() *notify.Config {
	return &b.appConfig
}

func (b *Builder) Context() context.Context {
	return b.ctx
}

func (b *Builder) UseAppConfig(cfg notify.Config) {
	b.appConfig = cfg
}

func (b *Builder) UseContext(ctx context.Context) {
	b.ctx = ctx
}

func (b *Builder) IocContainer() *ioc.Container {
	return b.ioc
}

func New(container *ioc.Container) *Builder {
	return &Builder{
		appConfig: notify.Config{ID: uuid.New().String()},
		ioc:       container,
		ctx: 	   context.Background(),
	}
}