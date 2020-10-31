package dnslb

import (
	"errors"
	"github.com/RomanIschenko/notify"
	"io"
)

type Registration struct {
	b notify.Broker
	pingHandler io.Closer
	app *notify.App
}

func (r Registration) Close() error {
	r.b.Emit(notify.BrokerEvent{
		AppID:    r.app.ID(),
		Event:    notify.BrokerAppDownEvent,
	})
	return r.pingHandler.Close()
}

func Register(app *notify.App, broker notify.Broker, addr string) (io.Closer, error) {
	if app == nil {
		return nil, errors.New("can not register nil app")
	}

	if broker == nil {
		return nil, errors.New("can not register app in nil broker")
	}

	broker.Emit(notify.BrokerEvent{
		Data:     addr,
		AppID:    app.ID(),
		Event:    notify.BrokerAppUpEvent,
	})

	pingHandler := handlePingEvent(broker, func() string {
		return addr
	})

	return Registration{broker, pingHandler, app}, nil
}