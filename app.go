package swirl

import (
	"context"
	"github.com/ischenkx/swirl/internal/emitter"
	"github.com/ischenkx/swirl/internal/pubsub"
	authmock "github.com/ischenkx/swirl/pkg/auth/mock"
	evcodec "github.com/ischenkx/swirl/pkg/default/event_codec"
	"reflect"
	"time"
)

type Config struct {
	CleanInterval time.Duration
	PubSubConfig  pubsub.Config
	EventsCodec   emitter.EventsCodec
	Auth          Authorizer
	Adapter       Adapter
}

func (c *Config) validate() {
	if c.Auth == nil {
		c.Auth = authmock.New()
	}
	if c.CleanInterval <= 0 {
		c.CleanInterval = time.Second * 15
	}
	if c.EventsCodec == nil {
		c.EventsCodec = evcodec.JSON{}
	}
}

type App struct {
	emitter          *emitter.Emitter
	pubsub           *pubsub.PubSub
	events           *eventsRegistry
	metricsCollector *metricsCollector
	adapter          Adapter
	auth             Authorizer
	config           Config
}

func (app *App) localMetrics() Metrics {
	return app.metricsCollector.get()
}

func (app *App) Metrics(flags ...interface{}) Metrics {
	if app.adapter == nil {
		return app.localMetrics()
	}
	for _, flag := range flags {
		if _, ok := flag.(LocalFlag); ok {
			return app.localMetrics()
		}
	}
	return app.adapter.Metrics()
}

func (app *App) Adapter() Adapter {
	return app.adapter
}

// Events are useful for node-bound event handling
func (app *App) Events(priority Priority) AppEvents {
	return app.events.forApp(priority)
}

func (app *App) localClient(id string) Client {
	return localClient{
		id:  id,
		app: app,
	}
}

func (app *App) localUser(id string) User {
	return localUser{
		id:  id,
		app: app,
	}
}

func (app *App) localTopic(id string) Topic {
	return localTopic{
		app: app,
		id:  id,
	}
}

func (app *App) Client(id string, flags ...interface{}) Client {
	for _, f := range flags {
		switch f.(type) {
		case LocalFlag:
			return app.localClient(id)
		}
	}
	if app.adapter != nil {
		return app.adapter.Client(id)
	}
	return app.localClient(id)
}

func (app *App) User(id string, flags ...interface{}) User {
	for _, f := range flags {
		switch f.(type) {
		case LocalFlag:
			return app.localUser(id)
		}
	}
	if app.adapter != nil {
		return app.adapter.User(id)
	}
	return app.localUser(id)
}

func (app *App) Topic(id string, flags ...interface{}) Topic {
	for _, f := range flags {
		switch f.(type) {
		case LocalFlag:
			return app.localTopic(id)
		}
	}
	if app.adapter != nil {
		return app.adapter.Topic(id)
	}
	return app.localTopic(id)
}

func (app *App) startCleaner(ctx context.Context) {
	ticker := time.NewTicker(app.config.CleanInterval)
	for {
		select {
		case <-ticker.C:
			ts := time.Now().UnixNano()
			cl := ChangeLog{TimeStamp: ts}
			clients := app.pubsub.Clean()
			cl.ClientsDown = append(cl.ClientsDown, clients...)
			for _, id := range clients {
				c := app.pubsub.Users().Delete(app.pubsub.Users().UserByClient(id), id, ts)
				cl.Merge(c)
			}
			app.events.callChange(cl)
		case <-ctx.Done():
			return
		}
	}
}

func (app *App) Start(ctx context.Context) {
	if app.adapter != nil {
		go app.adapter.Start(ctx)
	}
	go app.startCleaner(ctx)
	go app.pubsub.Start(ctx)
}

func (app *App) Server() Server {
	return Server{app: app}
}

func (app *App) On(event string, handler interface{}) {
	if err := app.emitter.Handle(event, handler); err != nil {
		app.events.callError(HandlerInitializationError{
			Reason:    err,
			EventName: event,
			Handler:   handler,
		})
	}
}

func (app *App) initEmitter() {
	var clientExample Client
	em := emitter.New(app.config.EventsCodec, reflect.TypeOf(&App{}), reflect.TypeOf(&clientExample).Elem())
	app.emitter = em
}

func (app *App) initEventsRegistry() {
	app.events = newEventsRegistry(app)
}

func (app *App) initMetricsCollector() {
	app.metricsCollector = newMetricsCollector(app)
}

func (app *App) initAdapter() {
	if app.adapter == nil {
		return
	}
	app.adapter.Init(app)
}

func New(config Config) *App {
	config.validate()
	app := &App{
		pubsub:  pubsub.New(config.PubSubConfig),
		adapter: config.Adapter,
		auth:    config.Auth,
		config: config,
	}
	app.initAdapter()
	app.initEventsRegistry()
	app.initEmitter()
	app.initMetricsCollector()
	return app
}
