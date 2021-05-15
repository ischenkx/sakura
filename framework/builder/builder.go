package builder

import (
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/framework/ioc"
)

const ImportPath = "github.com/RomanIschenko/notify/framework/builder"

type LabelOption struct {
	Name string
}

func WithLabel(l string) LabelOption {
	return LabelOption{Name: l}
}

type Builder struct {
	appConfig notify.Config
	ioc *ioc.Container
}

func (b *Builder) AppConfig() *notify.Config {
	return &b.appConfig
}

func (b *Builder) SetAppConfig(cfg notify.Config) {
	b.appConfig = cfg
}

func (b *Builder) Inject(data interface{}, opts ...interface{}) {

	var label string

	for _, opt := range opts {
		switch o := opt.(type) {
		case LabelOption:
			label = o.Name
		}
	}

	b.ioc.Add(ioc.NewEntry(data, label))
}

func New(container *ioc.Container) *Builder {
	return &Builder{
		appConfig: notify.Config{},
		ioc:       container,
	}
}