package swirl

import "context"

type Handle interface {
	Close()
}

type AdapterBuilder interface {
	Build(app *App) Adapter
}

type Adapter interface {
	Client(string) Client
	User(string) User
	Topic(string) Topic
	Emit(EventOptions)
	Metrics() Metrics
	Handle(string, interface{}) (Handle, error)
	Start(context.Context)
}

