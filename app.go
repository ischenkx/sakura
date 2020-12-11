package notify

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/RomanIschenko/notify/pubsub/namespace"
	"github.com/sirupsen/logrus"
	"runtime"
)

var logger = logrus.WithField("source", "notify_app")

type ServerConfig struct {
	Server  Server
	Workers int
	// if data handler returns non-nil error then client will be disconnected
	DataHandler func(*App, IncomingData)
}

func (cfg *ServerConfig) validate() {
	if cfg.Workers <= 0 {
		cfg.Workers = runtime.NumCPU()
	}
}

type Config struct {
	ID           string
	PubSubConfig pubsub.Config
	ServerConfig ServerConfig
	Auth         Auth
}
func (cfg *Config) validate() {
	if cfg.ID == "" {
		panic("cannot create app with empty id")
	}
}

type App struct {
	id           string
	pubsub       *pubsub.Pubsub
	auth         Auth
	serverConfig ServerConfig
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
	return client, err
}

func (app *App) inactivateClient(client *pubsub.Client) {
	if client == nil {
		return
	}
	app.pubsub.InactivateClient(client)
}

func (app *App) handle(data IncomingData) {
	if data.Client == nil {
		return
	}
	if data.Payload == nil {
		return
	}
	if app.serverConfig.DataHandler == nil {
		data.Payload = data.Payload[:0]
		return
	}
	app.serverConfig.DataHandler(app, data)
}

func (app *App) startServer(ctx context.Context) {
	if app.serverConfig.Server == nil {
		return
	}
	app.serverConfig.validate()
	server := app.serverConfig.Server
	for i := 0; i < app.serverConfig.Workers; i++ {
		go func(ctx context.Context, server Server) {
			for {
				select {
				case <-ctx.Done():
					return
				case conn := <-server.Accept():
					go func() {
						var resolved ResolvedConnection
						resolved.Client, resolved.Err = app.connect(conn.Opts, conn.AuthData)
						conn.Resolver <- resolved
					}()
				}
			}
		}(ctx, server)
	}

	for i := 0; i < app.serverConfig.Workers; i++ {
		go func(ctx context.Context, server Server) {
			for {
				select {
				case <-ctx.Done():
					return
				case client := <-server.Inactive():
					app.inactivateClient(client)
				}
			}
		}(ctx, server)
	}

	for i := 0; i < app.serverConfig.Workers; i++ {
		go func(ctx context.Context, server Server) {
			for {
				select {
				case <-ctx.Done():
					return
				case opts := <-server.Incoming():
					go app.handle(opts)
				}
			}
		}(ctx, server)
	}
}

func (app *App) Clients() []string {
	return app.pubsub.Clients()
}

func (app *App) Users() []string {
	return app.pubsub.Users()
}

func (app *App) Topics() []string {
	return app.pubsub.Topics()
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Proxy(ctx context.Context) *pubsub.Proxy {
	return app.pubsub.Proxy(ctx)
}

func (app *App) Events(ctx context.Context) *pubsub.EventsHub {
	return app.pubsub.Events(ctx)
}

func (app *App) NamespaceRegistry() *namespace.Registry {
	return app.pubsub.NamespaceRegistry()
}

func (app *App) Metrics() pubsub.Metrics {
	return app.pubsub.Metrics()
}

func (app *App) Publish(opts pubsub.PublishOptions) {
	err := app.pubsub.Publish(opts)
	if err != nil {
		logger.Debug("failed to publish:", err)
		return
	}
}

func (app *App) Subscribe(opts pubsub.SubscribeOptions) changelog.Log {
	return app.pubsub.Subscribe(opts)
}

func (app *App) Unsubscribe(opts pubsub.UnsubscribeOptions) changelog.Log {
	return app.pubsub.Unsubscribe(opts)
}

func (app *App) Disconnect(opts pubsub.DisconnectOptions) changelog.Log {
	return app.pubsub.Disconnect(opts)
}

func (app *App) Start(ctx context.Context) {
	logger.Infof("%s has started", app.id)
	app.startServer(ctx)
	go app.pubsub.StartCleaner(ctx)
}

func New(config Config) *App {
	config.validate()
	app := &App{
		id:           config.ID,
		pubsub:       pubsub.New(config.PubSubConfig),
		auth:         config.Auth,
		serverConfig: config.ServerConfig,
	}
	return app
}