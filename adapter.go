package swirl

import "context"

type Handle interface {
	Close()
}

type Adapter interface {
	Client(*App, string) Client
	User(*App, string) User
	Topic(*App, string) Topic
	Emit(EventOptions)
	Metrics() Metrics
	Handle(string, interface{}) (Handle, error)
	Start(context.Context, *App)
	Init(app *App)
}

