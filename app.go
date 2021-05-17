package notify

import (
	"context"
	"github.com/ischenkx/notify/internal/emitter"
	"github.com/ischenkx/notify/internal/events"
	"github.com/ischenkx/notify/internal/pubsub"
	"github.com/ischenkx/notify/internal/pubsub/message"
	"github.com/ischenkx/notify/pkg/default/event_codec"
	"log"
	"reflect"
	"time"
)

type Config struct {
	ID     string
	Auth   Auth
	PubSub pubsub.Config
	Codec  EventsCodec
}

func (cfg *Config) validate() {
	if cfg.ID == "" {
		panic("cannot create app with empty id")
	}

	if cfg.PubSub.CleanInterval <= 0 {
		cfg.PubSub.CleanInterval = time.Minute
	}

	if cfg.Codec == nil {
		cfg.Codec = evcodec.JSON{}
	}
}

type App struct {
	id      string
	pubsub  *pubsub.PubSub
	hooks   *events.Source
	emitter *emitter.Emitter
	auth    Auth
	config  Config
}

// calls a specified hook, additionally passes *App as an argument
func (app *App) callHook(name string, args ...interface{}) {
	app.hooks.Emit(name, append(args, app)...)
}

func (app *App) inactivate(id string) {
	client, changelog, err := app.pubsub.Inactivate(id, time.Now().UnixNano())
	if err != nil {
		return
	}
	app.callHook(changeHookName, changelog)
	app.callHook(inactivateHookName, Client{
		id:        client.ID(),
		app:       app,
		rawClient: client,
	})
}

func (app *App) handleMessage(id string, data []byte) {
	name, payload, err := app.emitter.DecodeRawData(data)
	if err != nil {
		log.Println(err)
		return
	}
	hnd, ok := app.emitter.GetHandler(name)
	if ok {
		d, ok := hnd.GetData(payload)
		if !ok {
			app.callHook(failedReceiveHookName, app.Client(id), name, payload)
		} else {
			app.callHook(receiveHookName, app.Client(id), IncomingEvent{
				Name: name,
				Data: d,
			})
			hnd.Call(app, app.Client(id), d)
		}
	}
}

func (app *App) startCleaner(ctx context.Context) {
	ticker := time.NewTicker(app.config.PubSub.CleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			clients, changelog := app.pubsub.Clean()
			for _, c := range clients {
				app.callHook(disconnectHookName, Client{
					id:        c.ID(),
					app:       app,
					rawClient: c,
				})
			}
			app.callHook(changeHookName, changelog)
		}
	}
}

func (app *App) ID() string {
	return app.id
}

func (app *App) SubscribeClient(opts SubscribeClientOptions) error {
	app.callHook(beforeClientSubscribeHookName, &opts)
	changelog, err := app.pubsub.SubscribeClient(opts)
	app.callHook(changeHookName, changelog)
	if err == nil {
		app.callHook(clientSubscribeHookName, app.Client(opts.ID), opts)
	}
	return err
}

func (app *App) UnsubscribeClient(opts UnsubscribeClientOptions) error {
	app.callHook(beforeClientUnsubscribeHookName, &opts)
	changelog, err := app.pubsub.UnsubscribeClient(opts)
	app.callHook(changeHookName, changelog)
	if err == nil {
		app.callHook(clientUnsubscribeHookName, app.Client(opts.ID), opts)
	}
	return err
}

func (app *App) SubscribeUser(opts SubscribeUserOptions) error {
	app.callHook(beforeUserSubscribeHookName, &opts)
	changelog, err := app.pubsub.SubscribeUser(opts)
	app.callHook(changeHookName, changelog)
	if err == nil {
		app.callHook(userSubscribeHookName, opts)
	}
	return err
}

func (app *App) UnsubscribeUser(opts UnsubscribeUserOptions) error {
	app.callHook(beforeUserUnsubscribeHookName, &opts)
	changelog, err := app.pubsub.UnsubscribeUser(opts)
	app.callHook(changeHookName, changelog)
	if err == nil {
		app.callHook(userUnsubscribeHookName, opts)
	}
	return err
}

func (app *App) Disconnect(opts DisconnectOptions) {
	app.callHook(beforeDisconnectHookName, &opts)
	clients, changelog := app.pubsub.Disconnect(opts)
	app.callHook(changeHookName, changelog)
	for _, c := range clients {
		app.callHook(disconnectHookName, c)
	}
}

func (app *App) Start(ctx context.Context) {
	app.pubsub.Start(ctx)
	go app.startCleaner(ctx)
}

func (app *App) Server() Server {
	return (*server)(app)
}

func (app *App) Hooks(p Priority) Hooks {
	return Hooks{
		app.hooks.NewHub(p),
	}
}

func (app *App) Client(id string) Client {
	return Client{
		id:        id,
		app:       app,
		rawClient: nil,
	}
}

func (app *App) User(id string) User {
	return User{
		id:  id,
		app: app,
	}
}

func (app *App) connect(opts ConnectOptions, auth string) (string, error) {
	if app.auth != nil {
		id, user, err := app.auth.Authorize(auth)
		if err != nil {
			return "", err
		}
		opts.ClientID = id
		opts.UserID = user
	}

	app.callHook(beforeConnectHookName, &opts)
	res := app.pubsub.Connect(opts)

	if res.Error != nil {
		return "", res.Error
	}

	c := Client{
		id:        res.Client.ID(),
		app:       app,
		rawClient: res.Client,
	}

	app.callHook(changeHookName, res.ChangeLog)
	if res.Reconnected {
		app.callHook(reconnectHookName, opts, c)
	} else {
		app.callHook(connectHookName, opts, c)
	}

	return c.id, nil
}

func (app *App) On(event string, handler interface{}) {
	app.emitter.Handle(event, handler)
}

func (app *App) Off(event string) {
	app.emitter.DeleteHandler(event)
}

func (app *App) Emit(ev Event) {
	app.callHook(beforeEmitHookName, &ev)
	data, err := app.emitter.EncodeRawData(ev.Name, ev.Data)
	if err != nil {
		log.Println("failed to emit:", err)
		return
	}

	app.pubsub.Publish(pubsub.PublishOptions{
		Clients:   ev.Clients,
		Users:     ev.Users,
		Topics:    ev.Topics,
		Data:      []message.Message{message.New(data)},
		TimeStamp: ev.Timestamp,
		Meta:      ev.Meta,
	})

	app.callHook(emitHookName, ev)
	app.callHook(rawEmitHookName, rawEvent{
		Payload: data,
		Meta:    ev.Meta,
		Clients: ev.Clients,
		Users:   ev.Users,
		Topics:  ev.Topics,
	})
}

func New(config Config) *App {
	config.validate()

	app := &App{
		id:      config.ID,
		pubsub:  pubsub.New(config.PubSub),
		hooks:   events.New(),
		auth:    config.Auth,
		config:  config,
		emitter: emitter.New(config.Codec, reflect.TypeOf(&App{}), reflect.TypeOf(Client{})),
	}

	return app
}