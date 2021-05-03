package notify

import (
	"context"
	"github.com/RomanIschenko/notify/default/batchproto"
	evcodec "github.com/RomanIschenko/notify/default/event_codec"
	basicps "github.com/RomanIschenko/notify/default/pubsub"
	"github.com/RomanIschenko/notify/internal/events"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/message"
	"github.com/RomanIschenko/notify/pubsub/protocol"
	"log"
	"time"
)

type Config struct {
	ID            string
	Auth          Auth
	PubSub        pubsub.Engine
	Protocol 	  protocol.Provider
	CleanInterval time.Duration
	Codec 		  EventsCodec
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

	if cfg.Codec == nil {
		cfg.Codec = evcodec.JSON{}
	}
}

type App struct {
	id     string
	pubsub pubsub.Engine
	hooks *events.Source
	emitter *emitter
	auth   Auth
	config Config
}

func (app *App) inactivate(id string) {
	client, changelog, err := app.pubsub.Inactivate(id, time.Now().UnixNano())
	if err != nil {
		return
	}
	app.hooks.Emit(changeHookName, app, changelog)
	app.hooks.Emit(inactivateHookName, app,  Client{client})
}

func (app *App) handleMessage(data []byte, client Client) {
	name, payload, err := app.emitter.decodeData(data)
	if err != nil {
		log.Println()
	}
	hnd, ok := app.emitter.getHandler(name)
	if ok {
		d, _ := hnd.decodeData(payload)
		app.hooks.Emit(receiveHookName, app,  Client{client}, IncomingEvent{
			Name: name,
			Data: d,
		})
		hnd.call(app, client, d)
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
				app.hooks.Emit(disconnectHookName, app,  Client{c})
			}
			app.hooks.Emit(changeHookName, app, changelog)
		}
	}
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Subscribe(opts SubscribeOptions) SubscriptionAlterationResult {
	app.hooks.Emit(beforeSubscribeHookName, app, &opts)
	changelog := app.pubsub.Subscribe(opts)
	app.hooks.Emit(changeHookName, app, changelog)
	return SubscriptionAlterationResult{changelog}
}

func (app *App) Unsubscribe(opts UnsubscribeOptions) SubscriptionAlterationResult {
	app.hooks.Emit(beforeUnsubscribeHookName, app, &opts)
	changelog := app.pubsub.Unsubscribe(opts)
	app.hooks.Emit(changeHookName, app, changelog)
	return SubscriptionAlterationResult{changelog}
}

func (app *App) Disconnect(opts DisconnectOptions) {
	app.hooks.Emit(beforeDisconnectHookName, app, &opts)

	clients, changelog := app.pubsub.Disconnect(opts)

	app.hooks.Emit(changeHookName, app, changelog)

	for _, c := range clients {
		app.hooks.Emit(disconnectHookName, app, c)
	}
}

func (app *App) Start(ctx context.Context) {
	app.pubsub.Start(ctx)
	go app.startCleaner(ctx)
}

func (app *App) Servable() Servable {
	return (*servable)(app)
}

func (app *App) Hooks(p Priority) Hooks {
	return Hooks {
		app.hooks.NewHub(p),
	}
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
			return Client{}, err
		}
		opts.ClientID = id
		opts.UserID = user
	}

	app.hooks.Emit(beforeConnectHookName, app, &opts)
	client, changelog, reconnected, err := app.pubsub.Connect(opts)

	if err != nil {
		return Client{}, err
	}

	app.hooks.Emit(changeHookName, app, changelog)
	if reconnected {
		log.Println("reconnected")
		app.hooks.Emit(reconnectHookName, app, opts, Client{client})
	} else {
		app.hooks.Emit(connectHookName, app, opts, Client{client})
	}

	return  Client{client}, nil
}

func (app *App) Action() ActionBuilder {
	return ActionBuilder{app: app}
}

func (app *App) On(event string, handler interface{}) {
	app.emitter.handle(event, handler)
}

func (app *App) Off(event string) {
	app.emitter.deleteHandler(event)
}

func (app *App) Emit(ev Event) {
	app.hooks.Emit(beforeEmitHookName, app, &ev)
	data, err := app.emitter.encodeData(ev.Name, ev.Data)
	if err != nil {
		log.Println("failed to emit:", err)
		return
	}

	app.pubsub.Publish(pubsub.PublishOptions{
		Clients:   ev.Clients,
		Users:     ev.Users,
		Topics:    ev.Topics,
		Data:     []message.Message{message.New(data)},
		TimeStamp: ev.Timestamp,
		Meta:      ev.Meta,
	})
}

func New(config Config) *App {
	config.validate()
	app := &App{
		id:      config.ID,
		pubsub:  config.PubSub,
		hooks:   events.New(),
		auth:    config.Auth,
		config:  config,
	}

	app.emitter = &emitter{
		app:      app,
		handlers: map[string]*eventHandler{},
		codec:    config.Codec,
	}
	return app
}
