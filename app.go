package notify

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"time"
)

var (
	IncorrectTransportErr = errors.New("transport is nil")
)

type AppConfig struct {
	ID string
	Messages MessageStorage
	Broker Broker
	PubSub PubSubConfig
	Auth Auth
	CleanInterval time.Duration
}

type App struct {
	id string
	events *EventsHandler
	pubsub *PubSub
	messages MessageStorage
}

type AppServer struct {
	auth Auth
	broker Broker
	cleanInterval time.Duration
	starter chan struct{}
	*App
}

func (app *App) generateClientID() string {
	return fmt.Sprintf("%s-%s", app.id, uuid.New().String())
}

func (app *App) identifyMessage(opts MessageOptions) SendOptions {
	mes := Message{
		Data: opts.Data,
		ID:   fmt.Sprintf("%s-%s", app.id, uuid.New().String()),
	}

	return SendOptions{
		opts.Users,
		opts.Clients,
		opts.Channels,
		mes,
	}
}

func (app *App) ID() string {
	return app.id
}

func (app *App) On() *EventsHandler {
	return app.events
}

func (app *App) send(opts SendOptions) {
	go app.pubsub.Send(opts)
}

func (app *App) join(opts JoinOptions) {
	app.pubsub.Join(opts)
}

func (app *App) leave(opts LeaveOptions) {
	app.pubsub.Leave(opts)
}

func (app *App) loadMessages(ids []string) ([]Message, []string) {
	if app.messages == nil {
		return nil, nil
	}
	return app.messages.Load(app.id, ids)
}

// highly not recommended to call it
// because it's automatically called at a specified interval
// in node or in app
func (app *App) Clean(ctx context.Context) {
	start := time.Now()
	ids := app.pubsub.Clean()
	elapsed := time.Since(start)
	log.Printf("Cleaning:\n\tDeleted: %v\n\tTime: %v\n", len(ids), elapsed)
	if app.messages != nil {
		app.messages.Clean()
	}
}

func (app *App) Send(mes MessageOptions) {
	opts := app.identifyMessage(mes)
	app.send(opts)
	app.events.handle(opts)
	if app.messages != nil {
		app.messages.Store(app.id, opts.Message)
	}
}

func (app *App) Join(opts JoinOptions) {
	app.join(opts)
	app.events.handle(opts)
}

func (app *App) Leave(opts LeaveOptions) {
	app.leave(opts)
	app.events.handle(opts)
}

func (app *App) connect(info ClientInfo, transport Transport) (*Client, error) {
	if transport == nil {
		return nil, IncorrectTransportErr
	}
	if info.ID == NilId {
		info.ID = app.generateClientID()
	}
	return app.pubsub.Connect(info, transport)
}

func (app *App) Disconnect(clientId string) {
	app.pubsub.Disconnect(clientId)
}

func NewAppServer(config AppConfig) Server {
	app := &App{
		id:       config.ID,
		events:   &EventsHandler{},
		messages: nil,
	}
	app.pubsub = newPubsub(app, config.PubSub)

	server := &AppServer{
		auth:          config.Auth,
		broker:        config.Broker,
		cleanInterval: config.CleanInterval,
		App:           app,
		starter: 	   make(chan struct{}),
	}

	return server
}



func (app *AppServer) DisconnectClient(client *Client) {
	app.pubsub.DisconnectClient(client)
}

//Connect takes two arguments: data and transport.
//If ServerApp has Auth then data is passed to Auth.Verify
//and gets client's info from there, else Connect will try to
//convert data to ClientInfo.
func (app *AppServer) Connect(data interface{}, transport Transport) (*Client, error) {
	if transport == nil || data == nil {
		return nil, errors.New("invalid data or transport")
	}

	var (
		clientInfo ClientInfo
		ok bool
		)

	if app.auth != nil {
		clientInfo, ok = app.auth.Verify(data)
	} else {
		clientInfo, ok = data.(ClientInfo)
	}
	if !ok {
		return nil, errors.New("failed to get client info")
	}
	if clientInfo.AppID != app.id {
		return nil, errors.New("wrong app")
	}
	return app.connect(clientInfo, transport)
}

func (app *AppServer) Handle(client *Client, r io.Reader) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if r == nil {
		return nil
	}
	err := app.events.handleData(client, r)
	if err != nil {
		app.pubsub.DisconnectClient(client)
	}
	return err
}

func (app *AppServer) Run(ctx context.Context) {
	app.starter <- struct{}{}
	defer func() {
		<-app.starter
	}()

	cleaner := time.NewTicker(app.cleanInterval)
	if app.broker != nil {
		app.broker.Handle(func(action BrokerAction) error {
			select {
			case <-ctx.Done():
				return errors.New("ctx is done")
			default:
			}

			if opts, ok := action.Data.(SendOptions); ok {
				app.send(opts)
			}
			if opts, ok := action.Data.(JoinOptions); ok {
				app.join(opts)
			}
			if opts, ok := action.Data.(LeaveOptions); ok {
				app.leave(opts)
			}
			return nil
		})

		app.broker.Emit(BrokerAction{
			Event: InstanceUpEvent,
		})
		app.broker.Emit(BrokerAction{
			AppID: app.id,
			Event: AppUpEvent,
		})
		defer app.broker.Emit(BrokerAction{
			Event: InstanceDownEvent,
		})
		defer app.broker.Emit(BrokerAction{
			AppID: app.id,
			Event: AppDownEvent,
		})
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-cleaner.C:
			go app.Clean(ctx)
		}
	}
}