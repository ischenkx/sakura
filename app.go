package notify

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var logger = logrus.WithField("source", "notify_app")

type ServerConfig struct {
	Instance  Server
	Workers int
}

func (cfg *ServerConfig) validate() {
	if cfg.Workers <= 0 {
		cfg.Workers = runtime.NumCPU()
	}
}

type Config struct {
	ID           string
	PubsubConfig pubsub.Config
	CleanInterval time.Duration
	Server ServerConfig
	Auth         Auth
}

func (cfg *Config) validate() {
	if cfg.ID == "" {
		panic("cannot create app with empty id")
	}
}

type App struct {
	id           string
	hub 		 *pubsub.Hub
	auth         Auth
	config Config
	proxyRegistry *proxyRegistry
	eventsRegistry *eventsRegistry
}

func (app *App) Proxy(ctx context.Context) *Proxy {
	return app.proxyRegistry.hub(ctx)
}

func (app *App) connect(opts pubsub.ConnectOptions, auth string) (pubsub.Client, error) {
	if app.auth != nil {
		clientID, userID, err := app.auth.Authorize(auth)
		if err != nil {
			logger.Debug("failed to authorize:", err)
			return pubsub.Client{}, err
		}
		opts.ClientID = clientID
		opts.UserID = userID
	}
	app.proxyRegistry.emitConnect(app, &opts)
	c, log, err := app.hub.Connect(opts)
	if err != nil {
		logger.Debug("failed to connect:", err)
	} else {
		app.eventsRegistry.emitConnect(app, opts, c, log)
	}

	return c, err
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Events(ctx context.Context) *EventsHub {
	return app.eventsRegistry.hub(ctx)
}

func (app *App) Publish(opts pubsub.PublishOptions) {
	app.proxyRegistry.emitPublish(app, &opts)
	app.hub.Publish(opts)
	app.eventsRegistry.emitPublish(app, opts)
}

func (app *App) Subscribe(opts pubsub.SubscribeOptions) {
	app.proxyRegistry.emitSubscribe(app, &opts)
	log := app.hub.Subscribe(opts)
	app.eventsRegistry.emitSubscribe(app, opts, log)
}

func (app *App) Unsubscribe(opts pubsub.UnsubscribeOptions) {
	app.proxyRegistry.emitUnsubscribe(app, &opts)
	log := app.hub.Unsubscribe(opts)
	app.eventsRegistry.emitUnsubscribe(app, opts, log)
}

func (app *App) Disconnect(opts pubsub.DisconnectOptions) {
	app.proxyRegistry.emitDisconnect(app, &opts)
	log := app.hub.Disconnect(opts)
	app.eventsRegistry.emitDisconnect(app, opts, log)
}

func (app *App) startServer(ctx context.Context) {
	if app.config.Server.Instance != nil {
		go app.config.Server.Instance.Start(ctx, (*servable)(app))
	}
}

func (app *App) startHubCleaner(ctx context.Context) {
	ticker := time.NewTicker(app.config.CleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log := app.hub.Clean()
			app.eventsRegistry.emitChange(app, log)
		}
	}
}

func (app *App) Start(ctx context.Context) {
	logger.Infof("%s has started", app.id)
	app.startServer(ctx)
	app.hub.Start(ctx)
}

func New(config Config) *App {
	config.validate()
	h := pubsub.New(config.PubsubConfig)
	app := &App{
		id:  config.ID,
		hub: h,
		auth: config.Auth,
		eventsRegistry: newEventsRegistry(),
		proxyRegistry: newProxyRegistry(),
		config: config,
	}
	return app
}