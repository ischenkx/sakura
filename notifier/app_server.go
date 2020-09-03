package notifier

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify"
	"io"
	"io/ioutil"
	"time"
)

type DataHandler func(client *notify.Client, r io.Reader) error

type AppServer struct {
	broker         	   Broker
	cleanInterval 	   time.Duration
	dataHandler   	   DataHandler
	starter     	   chan struct{}
	*notify.App
}

func (app *AppServer) Handle(client *notify.Client, r io.Reader) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if r == nil {
		return nil
	}
	if app.dataHandler == nil {
		ioutil.ReadAll(r)
		return nil
	}
	return app.dataHandler(client, r)
}

func (app *AppServer) runBrokerEventLoop(ctx context.Context) {
	if app.broker == nil {
		return
	}
	brokerSub := app.broker.Subscribe()
	defer brokerSub.Close()
	subscription := app.Events().Subscribe()
	defer subscription.Close()

	app.broker.Publish(BrokerMessage{
		Data:     ctx.Value("instanceUpArg"),
		AppID:    app.ID(),
		Event:    BrokerInstanceUpEvent,
	})

	defer app.broker.Publish(BrokerMessage{
		AppID:    app.ID(),
		Event:    BrokerInstanceDownEvent,
	})

	app.broker.Publish(BrokerMessage{
		AppID:    app.ID(),
		Event:    BrokerAppUpEvent,
	})

	defer app.broker.Publish(BrokerMessage{
		AppID:    app.ID(),
		Event:    BrokerAppDownEvent,
	})

	for {
		select {
		case appEvent := <-subscription.Channel():
			switch appEvent.Type {
			case notify.JoinEvent, notify.LeaveEvent, notify.SendEvent:
				app.broker.Publish(BrokerMessage{
					Data:     appEvent.Data,
					AppID:    app.ID(),
					Event:    appEvent.Type,
				})
			}
		case mes := <-brokerSub.Channel():
			if mes.AppID != app.ID() {
				continue
			}
			switch mes.Event {
			case notify.SendEvent:
				if opts, ok := mes.Data.(notify.SendOptions); ok {
					opts.Event = BrokerSendEvent
					go app.Send(opts)
				}
			case notify.JoinEvent:
				if opts, ok := mes.Data.(notify.JoinOptions); ok {
					opts.Event = BrokerJoinEvent
					go app.Join(opts)
				}
			case notify.LeaveEvent:
				if opts, ok := mes.Data.(notify.LeaveOptions); ok {
					opts.Event = BrokerLeaveEvent
					go app.Leave(opts)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (app *AppServer) Run(ctx context.Context) {
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

func NewAppServer(app *notify.App, config Config) *AppServer {
	return &AppServer{
		broker:        config.Broker,
		cleanInterval: config.CleanInterval,
		App:           app,
		dataHandler:   config.DataHandler,
		starter: 	   make(chan struct{}, 1),
	}
}
