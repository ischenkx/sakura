package notify

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify/events"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"runtime"
	"sync"
	"time"
)

const AppUpArgName = "appUpArg"

var logger = logrus.WithField("source", "notify_app")

type ServerConfig struct {
	Server Server
	Goroutines int
	// if data handler returns non-nil error then client will be disconnected
	DataHandler func(*App, IncomingData) error
}

func (cfg ServerConfig) validate() ServerConfig {
	if cfg.Goroutines <= 0 {
		cfg.Goroutines = runtime.NumCPU()
	}
	return cfg
}

type Config struct {
	ID       		string
	Broker	 		Broker
	PubSubConfig  	pubsub.Config
	ServerConfig 	ServerConfig
	Auth			Auth
}

func (cfg Config) validate() Config {
	if cfg.ID == "" {
		panic("cannot create app with empty id")
	}
	return cfg
}

type App struct {
	id          string
	events      *events.Source
	pubsub      *pubsub.Pubsub
	auth 		Auth
	serverConfig ServerConfig
	broker 	    Broker
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Events() *events.Source {
	return app.events
}

func (app *App) NSConfig(ns string) (pubsub.NamespaceConfig, bool) {
	return app.pubsub.NS().Get(ns)
}

func (app *App) RegisterNS(ns string, config pubsub.NamespaceConfig) {
	app.pubsub.NS().Register(ns, config)
}

func (app *App) UnregisterNS(ns string) {
	app.pubsub.NS().Unregister(ns)
}

func (app *App) Metrics() pubsub.Metrics {
	return app.pubsub.Metrics()
}

func (app *App) publish(opts pubsub.PublishOptions) pubsub.Result {
	return app.pubsub.Publish(opts)
}

func (app *App) subscribe(opts pubsub.SubscribeOptions) pubsub.Result {
	return app.pubsub.Subscribe(opts)
}

func (app *App) unsubscribe(opts pubsub.UnsubscribeOptions) pubsub.Result {
	return app.pubsub.Unsubscribe(opts)
}

func (app *App) connect(opts pubsub.ConnectOptions, auth string) (*pubsub.Client, error) {
	if app.auth != nil {
		clientID, err := app.auth.Authorize(auth)
		if err != nil {
			logger.Debug("failed to authorize:", err)
			return nil, err
		}
		opts.ID = clientID
	}
	client, err := app.pubsub.Connect(opts)
	if err != nil {
		logger.Debug("failed to connect:", err)
		return client, err
	}
	app.events.Emit(events.Event{
		Data: client,
		Type: ConnectEvent,
	})
	return client, err
}

func (app *App) inactivateClient(client *pubsub.Client) {
	if client == nil {
		return
	}
	app.pubsub.InactivateClient(client)
	app.events.Emit(events.Event{
		Data: client.ID(),
		Type: InactivateEvent,
	})
}

func (app *App) handle(data IncomingData) error {
	if data.Client == nil {
		return errors.New("client is nil")
	}
	if data.Reader == nil {
		return nil
	}
	if app.serverConfig.DataHandler == nil {
		_, err := ioutil.ReadAll(data.Reader)
		return err
	}
	return app.serverConfig.DataHandler(app, data)
}

func (app *App) startBrokerEventLoop(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

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
				app.broker.HandlePubsubResult(app.ID(), app.publish(opts))
			}
		case SubscribeEvent:
			if opts, ok := e.Data.(pubsub.SubscribeOptions); ok {
				app.broker.HandlePubsubResult(app.ID(), app.subscribe(opts))
			}
		case UnsubscribeEvent:
			if opts, ok := e.Data.(pubsub.UnsubscribeOptions); ok {
				app.broker.HandlePubsubResult(app.ID(), app.unsubscribe(opts))
			}
		}
	})
	defer brokerHandlerCloser.Close()

	appHandlerCloser := app.Events().Handle(func(e events.Event) {
		switch e.Type {
		case SubscribeEvent, UnsubscribeEvent, PublishEvent:
			if actionResult, ok := e.Data.(ActionResult); ok {
				input := actionResult.Input
				pubsubResult := actionResult.PubsubResult
				app.broker.HandlePubsubResult(app.ID(), pubsubResult)
				app.broker.Emit(BrokerEvent{
					Data:     input,
					AppID:    app.ID(),
					Event:    e.Type,
				})
			}
		}
	})
	defer appHandlerCloser.Close()

	<-ctx.Done()
}

func (app *App) startServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if app.serverConfig.Server == nil {
		return
	}
	server := app.serverConfig.Server
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-server.Inactive():
			app.inactivateClient(client)
		case conn := <-server.Accept():
			go func() {
				var resolved ResolvedConnection
				resolved.Client, resolved.Err = app.connect(conn.Opts, conn.AuthData)
				conn.Resolver <- resolved
			}()
		case opts := <-server.Incoming():
			app.handle(opts)
		}
	}
}

func (app *App) Publish(opts pubsub.PublishOptions) {
	res := app.pubsub.Publish(opts)
	app.events.Emit(events.Event{
		Data: ActionResult{
			Input:        opts,
			PubsubResult: res,
		},
		Type: PublishEvent,
	})
}

func (app *App) Subscribe(opts pubsub.SubscribeOptions) {
	app.events.Emit(events.Event{
		Data: ActionResult{
			Input:        opts,
			PubsubResult: app.subscribe(opts),
		},
		Type: SubscribeEvent,
	})
}

func (app *App) Unsubscribe(opts pubsub.UnsubscribeOptions) {
	app.events.Emit(events.Event{
		Data: ActionResult{
			Input:        opts,
			PubsubResult: app.unsubscribe(opts),
		},
		Type: UnsubscribeEvent,
	})
}

func (app *App) Disconnect(opts pubsub.DisconnectOptions) {
	app.events.Emit(events.Event{
		Data: ActionResult{
			Input:        opts,
			PubsubResult: app.pubsub.Disconnect(opts),
		},
		Type: DisconnectEvent,
	})
}

// starts app, returns a sync.WaitGroup to shutdown gracefully
func (app *App) Start(ctx context.Context) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logger.Debugf("app %s has started", app.id)

	app.events.Emit(events.Event{
		Data: time.Now(),
		Type: AppStartedEvent,
	})

	go app.startBrokerEventLoop(ctx, wg)
	for i := 0; i < app.serverConfig.Goroutines; i++ {
		wg.Add(1)
		go app.startServer(ctx, wg)
	}
	go app.pubsub.Start(ctx)

	//goroutine checking whether app is done
	go func() {
		select {
		case <-ctx.Done():
		}
		app.events.Emit(events.Event{
			Type: AppDoneEvent,
			Data: time.Now(),
		})
	}()

	return wg
}

func New(config Config) *App {
	config = config.validate()

	app := &App{
		id:       	 config.ID,
		events:   	 events.NewSource(),
		pubsub:   	 pubsub.New(config.PubSubConfig),
		auth:		 config.Auth,
		broker:      config.Broker,
		serverConfig: config.ServerConfig,
	}
	return app
}