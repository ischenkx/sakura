package notify

import (
	"context"
	"errors"
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
	brokerSub := app.broker.Subscribe()
	defer brokerSub.Close()
	subscription := app.Events().Subscribe()
	defer subscription.Close()

	app.broker.Publish(BrokerMessage{
		Data:  ctx.Value("instanceUpArg"),
		AppID: app.ID(),
		Event: BrokerInstanceUpEvent,
	})

	defer app.broker.Publish(BrokerMessage{
		AppID: app.ID(),
		Event: BrokerInstanceDownEvent,
	})

	app.broker.Publish(BrokerMessage{
		AppID: app.ID(),
		Event: BrokerAppUpEvent,
	})

	defer app.broker.Publish(BrokerMessage{
		AppID: app.ID(),
		Event: BrokerAppDownEvent,
	})

	for {
		select {
		case appEvent := <-subscription.Channel():
			switch appEvent.Type {
			case JoinEvent, LeaveEvent, SendEvent:
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
			case SendEvent:
				if opts, ok := mes.Data.(SendOptions); ok {
					opts.Event = BrokerSendEvent
					go app.Send(opts)
				}
			case JoinEvent:
				if opts, ok := mes.Data.(JoinOptions); ok {
					opts.Event = BrokerJoinEvent
					go app.Join(opts)
				}
			case LeaveEvent:
				if opts, ok := mes.Data.(LeaveOptions); ok {
					opts.Event = BrokerLeaveEvent
					go app.Leave(opts)
				}
			}
		case <-ctx.Done():
			return
		}
	}
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
