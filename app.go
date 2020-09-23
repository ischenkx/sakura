package notify

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify/events"
	"github.com/RomanIschenko/pubsub"
	"io"
	"io/ioutil"
)

type Config struct {
	ID       	string
	Broker	 	Broker
	PubSub  	pubsub.Config
	DataHandler func(client *pubsub.Client, r io.Reader) error
}

type App struct {
	id          string
	events      *events.Source
	pubsub      *pubsub.Pubsub
	dataHandler func(client *pubsub.Client, r io.Reader) error
	broker 	    Broker
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Events() *events.Source {
	return app.events
}

func (app *App) publish(opts pubsub.PublishOptions) {
	app.pubsub.Publish(opts)
}

func (app *App) subscribe(opts pubsub.SubscribeOptions) {
	app.pubsub.Subscribe(opts)
}

func (app *App) unsubscribe(opts pubsub.UnsubscribeOptions) {
	app.pubsub.Unsubscribe(opts)
}

func (app *App) Publish(opts pubsub.PublishOptions) {
	app.pubsub.Publish(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: PublishEvent,
	})
}

func (app *App) Subscribe(opts pubsub.SubscribeOptions) {
	app.subscribe(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: SubscribeEvent,
	})
}

func (app *App) Unsubscribe(opts pubsub.UnsubscribeOptions) {
	app.unsubscribe(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: UnsubscribeEvent,
	})
}

func (app *App) Connect(opts pubsub.ConnectOptions) (*pubsub.Client, error) {
	client, err := app.pubsub.Connect(opts)
	if err == nil {
		app.events.Emit(events.Event{
			Data: client,
			Type: ConnectEvent,
		})
	}
	return client, err
}

func (app *App) InactivateClient(client *pubsub.Client) {
	if client == nil {
		return
	}
	app.pubsub.InactivateClient(client)
	app.events.Emit(events.Event{
		Data: client.ID(),
		Type: InactivateEvent,
	})
}

func (app *App) Disconnect(opts pubsub.DisconnectOptions) {
	app.pubsub.Disconnect(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: DisconnectEvent,
	})
}

func (app *App) Handle(client *pubsub.Client, r io.Reader) error {
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

func (app *App) startBrokerEventLoop(ctx context.Context) {
	if app.broker == nil {
		return
	}
	brokerHandlerCloser := app.broker.Handle(func(e BrokerEvent) {
		if app == nil {
			return
		}
		if e.AppID != app.ID() {
			return
		}
		switch e.Event {
		case PublishEvent:
			if opts, ok := e.Data.(pubsub.PublishOptions); ok {
				app.publish(opts)
			}
		case SubscribeEvent:
			if opts, ok := e.Data.(pubsub.SubscribeOptions); ok {
				go app.subscribe(opts)
			}
		case UnsubscribeEvent:
			if opts, ok := e.Data.(pubsub.UnsubscribeOptions); ok {
				go app.unsubscribe(opts)
			}
		}
	})
	defer brokerHandlerCloser.Close()

	appHandlerCloser := app.Events().Handle(func(e events.Event) {
		switch e.Type {
		case SubscribeEvent, UnsubscribeEvent, PublishEvent:
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

func (app *App) Start(ctx context.Context) {
	go app.startBrokerEventLoop(ctx)
	app.pubsub.Start(ctx)
}

func New(config Config) *App {
	app := &App{
		id:       	 config.ID,
		events:   	 events.NewSource(),
		pubsub:   	 pubsub.New(config.PubSub),
		dataHandler: config.DataHandler,
		broker:      config.Broker,
	}

	return app
}