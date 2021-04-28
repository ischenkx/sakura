package notify

import (
	"context"
	"github.com/RomanIschenko/notify/default/batchproto"
	basicps "github.com/RomanIschenko/notify/default/pubsub"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/protocol"
	"time"
)

type Config struct {
	ID            string
	Auth          Auth
	PubSub        pubsub.Engine
	Protocol 	  protocol.Provider
	CleanInterval time.Duration
}

func (cfg *Config) validate() {
	if cfg.ID == "" {
		panic("cannot create app with empty id")
	}

	if cfg.CleanInterval <= 0 {
		cfg.CleanInterval = time.Minute
	}

	if cfg.PubSub == nil {
		cfg.PubSub = basicps.New(basicps.Config{
			ProtoProvider:    batchproto.NewProvider(1024),
			InvalidationTime: time.Minute,
			CleanInterval:    cfg.CleanInterval,
		})
	}
}

type App struct {
	id     string
	pubsub pubsub.Engine
	auth   Auth
	config Config
	proxy  *proxyRegistry
	events *eventsRegistry
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Subscribe(opts SubscribeOptions) ChangeLog {
	app.proxy.emitSubscribe(app, &opts)
	changelog := app.pubsub.Subscribe(opts)
	app.events.emitChange(app, changelog)
	return changelog
}

func (app *App) Unsubscribe(opts UnsubscribeOptions) ChangeLog {
	app.proxy.emitUnsubscribe(app, &opts)
	changelog := app.pubsub.Unsubscribe(opts)
	app.events.emitChange(app, changelog)
	return changelog
}

func (app *App) Publish(opts PublishOptions) {
	app.proxy.emitPublish(app, &opts)
	app.pubsub.Publish(opts)
}

func (app *App) Disconnect(opts DisconnectOptions) {
	app.proxy.emitDisconnect(app, &opts)
	clients, changelog := app.pubsub.Disconnect(opts)
	app.events.emitChange(app, changelog)
	for _, c := range clients {
		app.events.emitDisconnect(app, c)
	}
}

func (app *App) startCleaner(ctx context.Context) {
	ticker := time.NewTicker(app.config.CleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			clients, changelog := app.pubsub.Clean()
			for _, c := range clients {
				app.events.emitDisconnect(app, c)
			}
			app.events.emitChange(app, changelog)
		}
	}
}

func (app *App) Start(ctx context.Context) {
	app.pubsub.Start(ctx)
	go app.startCleaner(ctx)
}

func (app *App) Servable() Servable {
	return (*servable)(app)
}

func (app *App) Events(p Priority) *Events {
	return app.events.new(p)
}

func (app *App) Proxy(p Priority) *Proxy {
	return app.proxy.new(p)
}

func (app *App) IsSubscribed(client string, topic string) (bool, error) {
	return app.pubsub.IsSubscribed(client, topic)
}

func (app *App) IsUserSubscribed(user string, topic string) (bool, error) {
	return app.pubsub.IsUserSubscribed(user, topic)
}

func (app *App) TopicSubscribers(id string) ([]string, error) {
	return app.pubsub.TopicSubscribers(id)
}

func (app *App) ClientSubscriptions(id string) ([]string, error) {
	return app.pubsub.ClientSubscriptions(id)
}

func (app *App) UserSubscriptions(userID string) ([]string, error) {
	return app.pubsub.UserSubscriptions(userID)
}

func (app *App) connect(opts ConnectOptions, auth string) (Client, error) {
	if app.auth != nil {
		id, user, err := app.auth.Authorize(auth)
		if err != nil {
			return nil, err
		}
		opts.ClientID = id
		opts.UserID = user
	}

	app.proxy.emitConnect(app, &opts)
	client, changelog, reconnected, err := app.pubsub.Connect(opts)

	if err != nil {
		return nil, err
	}

	app.events.emitChange(app, changelog)
	if reconnected {
		app.events.emitReconnect(app, opts, client)
	} else {
		app.events.emitConnect(app, opts, client)
	}

	return client, nil
}

func (app *App) Action() ActionBuilder {
	return ActionBuilder{app: app}
}

func New(config Config) *App {
	config.validate()

	app := &App{
		id:     config.ID,
		pubsub: config.PubSub,
		auth:   config.Auth,
		events: newEventsRegistry(),
		proxy:  newProxyRegistry(),
		config: config,
	}
	return app
}
