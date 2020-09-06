package notify

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify/events"
	"io"
	"io/ioutil"
	"time"
)

type appServer struct {
	broker        Broker
	cleanInterval time.Duration
	dataHandler   func(client *Client, r io.Reader) error
	starter       chan struct{}
	*App
}

func (app *appServer) Handle(client *Client, r io.Reader) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if r == nil {
		return nil
	}
	if app.dataHandler == nil {
		_, err := ioutil.ReadAll(r)
		return err
	}
	return app.dataHandler(client, r)
}

func (app *appServer) runBrokerEventLoop(ctx context.Context) {
	if app.broker == nil {
		return
	}
	brokerHandlerCloser := app.broker.Handle(func(e BrokerEvent) {
		if e.AppID != app.ID() {
			return
		}
		switch e.Event {
		case SendEvent:
			if opts, ok := e.Data.(SendOptions); ok {
				opts.Event = BrokerSendEvent
				go app.Send(opts)
			}
		case JoinEvent:
			if opts, ok := e.Data.(JoinOptions); ok {
				opts.Event = BrokerJoinEvent
				go app.Join(opts)
			}
		case LeaveEvent:
			if opts, ok := e.Data.(LeaveOptions); ok {
				opts.Event = BrokerLeaveEvent
				go app.Leave(opts)
			}
		}
	})
	defer brokerHandlerCloser.Close()

	appHandlerCloser := app.Events().Handle(func(e events.Event) {
		switch e.Type {
		case JoinEvent, LeaveEvent, SendEvent:
			app.broker.Emit(BrokerEvent{
				Data:     e.Data,
				AppID:    app.ID(),
				Event:    e.Type,
			})
		}
	})
	defer appHandlerCloser.Close()

	app.broker.Emit(BrokerEvent{
		Data:  ctx.Value("instanceUpArg"),
		AppID: app.ID(),
		Event: BrokerInstanceUpEvent,
	})

	defer app.broker.Emit(BrokerEvent{
		AppID: app.ID(),
		Event: BrokerInstanceDownEvent,
	})

	app.broker.Emit(BrokerEvent{
		AppID: app.ID(),
		Event: BrokerAppUpEvent,
	})

	defer app.broker.Emit(BrokerEvent{
		AppID: app.ID(),
		Event: BrokerAppDownEvent,
	})

	<-ctx.Done()
}

func (app *appServer) Run(ctx context.Context) {
	app.starter <- struct{}{}
	defer func() {
		<-app.starter
	}()

	go app.runBrokerEventLoop(ctx)
	cleaner := time.NewTicker(app.cleanInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-cleaner.C:
			go app.Clean(ctx)
		}
	}
}
