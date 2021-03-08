package notify

import (
	"context"
	"github.com/RomanIschenko/notify/internal/pubsub"
	"time"
)

type Config struct {
	ID            string
	Auth          Auth
	CleanInterval time.Duration
	InvalidationTime time.Duration
}

func (cfg *Config) validate() {
	if cfg.ID == "" {
		panic("cannot create app with empty id")
	}

	if cfg.CleanInterval <= 0 {
		cfg.CleanInterval = time.Second * 25
	}

	if cfg.InvalidationTime < 0 {
		cfg.InvalidationTime = time.Minute
	}
}

type App struct {
	id             string
	pubsub         *pubsub.PubSub
	auth           Auth
	config         Config
	proxyRegistry  *proxyRegistry
	eventsRegistry *eventsRegistry
}

func (app *App) ID() string {
	return app.id
}

func(app *App) Subscribe(opts SubscribeOptions) {
	app.proxyRegistry.emitSubscribe(app, &opts)
	changelog := app.pubsub.Subscribe(opts)
	app.eventsRegistry.emitChange(app, changelog)
	// TODO
	// add unsubscribe event support
}

func(app *App) Unsubscribe(opts UnsubscribeOptions) {
	app.proxyRegistry.emitUnsubscribe(app, &opts)
	changelog := app.pubsub.Unsubscribe(opts)

	app.eventsRegistry.emitChange(app, changelog)

	// TODO
	// add unsubscribe event support
}

func(app *App) Publish(opts PublishOptions) {
	app.proxyRegistry.emitPublish(app, &opts)
	app.pubsub.Publish(opts)
}

func(app *App) Disconnect(opts DisconnectOptions)  {
	app.proxyRegistry.emitDisconnect(app, &opts)
	clients, changelog := app.pubsub.Disconnect(opts)
	app.eventsRegistry.emitChange(app, changelog)
	for _, c := range clients {
		app.eventsRegistry.emitDisconnect(app, c)
	}
}

func (app *App) startCleaner(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			clients, changelog := app.pubsub.Clean()
			for _, c := range clients {
				app.eventsRegistry.emitDisconnect(app, c)
			}
			app.eventsRegistry.emitChange(app, changelog)
		}
	}
}

func (app *App) Start(ctx context.Context) {
	app.pubsub.Start(ctx)
	go app.startCleaner(ctx)
}

func(app *App) Servable() Servable {
	return (*servable)(app)
}

func(app *App) Events() *AppEvents {
	return app.eventsRegistry.newHub()
}

func(app *App) Proxy() *Proxy {
	return app.proxyRegistry.newHub()
}

func(app *App) Metrics() Metrics {
	return Metrics{
		PubSubMetrics: app.pubsub.Metrics(),
	}
}

func(app *App) IsSubscribed(client string, topic string) (bool, error) {
	return app.pubsub.IsSubscribed(client, topic)
}

func(app *App) IsUserSubscribed(user string, topic string) (bool, error) {
	return app.pubsub.IsUserSubscribed(user, topic)

}

func(app *App) TopicSubscribers(id string) ([]string, error) {
	return app.pubsub.TopicSubscribers(id)
}

func(app *App) ClientSubscriptions(id string) ([]string, error) {
	return app.pubsub.ClientSubscriptions(id)

}

func(app *App) UserSubscriptions(userID string) ([]string, error) {
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

	app.proxyRegistry.emitConnect(app, &opts)
	client, changelog, reconnected, err := app.pubsub.Connect(opts)

	if err != nil {
		return nil, err
	}

	app.eventsRegistry.emitChange(app, changelog)
	if reconnected {
		app.eventsRegistry.emitReconnect(app, opts, client)
	} else {
		app.eventsRegistry.emitConnect(app, opts, client)
	}

	return client, nil
}

func New(config Config) *App {
	config.validate()
	h := pubsub.New(pubsub.Config{
		InvalidationTime: config.InvalidationTime,
		CleanInterval:    config.CleanInterval,
	})
	app := &App{
		id:             config.ID,
		pubsub:         h,
		auth:           config.Auth,
		eventsRegistry: newEventsRegistry(),
		proxyRegistry:  newProxyRegistry(),
		config:         config,
	}
	return app
}