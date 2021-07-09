package redaptor

import (
	"context"
	"github.com/ischenkx/swirl"
)

type Adapter struct {
	app *swirl.App
}

func (a *Adapter) Client(id string) swirl.Client {
	return newClient(a, a.app.Client(id, swirl.LocalFlag{}))
}

func (a *Adapter) User(s string) swirl.User {
	panic("implement me")
}

func (a *Adapter) Topic(s string) swirl.Topic {
	panic("implement me")
}

func (a *Adapter) Emit(options swirl.EventOptions) {
	panic("implement me")
}

func (a *Adapter) Metrics() swirl.Metrics {
	panic("implement me")
}

func (a *Adapter) Handle(s string, i interface{}) (swirl.Handle, error) {
	panic("implement me")
}

func (a Adapter) Start(ctx context.Context, app *swirl.App) {
	panic("implement me")
}

func New() *Adapter {
	return &Broker{}
}
