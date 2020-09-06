package notify

import (
	"context"
	"errors"
	"fmt"
	events "github.com/RomanIschenko/notify/events"
	"github.com/RomanIschenko/notify/message"
	"github.com/google/uuid"
	"log"
	"time"
)

var (
	IncorrectTransportErr = errors.New("transport is nil")
)

type App struct {
	id       string
	events   *events.Source
	pubsub   *PubSub
	messages message.Storage
}

func (app *App) generateClientID() string {
	return fmt.Sprintf("%s-%s", app.id, uuid.New().String())
}

func (app *App) identifyMessage(opts MessageSendOptions) SendOptions {
	mes := message.Message{
		Data: opts.Data,
		ID:   fmt.Sprintf("%s-%s", app.id, uuid.New().String()),
	}

	return SendOptions{
		Users:      opts.Users,
		Clients:    opts.Clients,
		Channels:   opts.Channels,
		ToBeStored: opts.ToBeStored,
		Message:    mes,
		Event:   	opts.Event,
	}
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Events() *events.Source {
	return app.events
}

func (app *App) send(opts SendOptions) {
	app.pubsub.Send(opts)
}

func (app *App) join(opts JoinOptions) {
	app.pubsub.Join(opts)
}

func (app *App) leave(opts LeaveOptions) {
	app.pubsub.Leave(opts)
}

func (app *App) loadMessages(ids []string) ([]message.Message, []string) {
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
}

func (app *App) Send(opts SendOptions) {
	if opts.Event == "" {
		opts.Event = SendEvent
	}
	app.pubsub.Send(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: opts.Event,
	})
	if app.messages != nil && opts.ToBeStored {
		app.messages.Store(app.id, opts.Message)
	}
}

func (app *App) SendMessage(mes MessageSendOptions) {
	if mes.Event == "" {
		mes.Event = SendEvent
	}
	opts := app.identifyMessage(mes)
	app.send(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: opts.Event,
	})
	if app.messages != nil && opts.ToBeStored {
		app.messages.Store(app.id, opts.Message)
	}
}

func (app *App) Join(opts JoinOptions) {
	if opts.Event == "" {
		opts.Event = JoinEvent
	}
	app.join(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: opts.Event,
	})
}

func (app *App) Leave(opts LeaveOptions) {
	if opts.Event == "" {
		opts.Event = LeaveEvent
	}
	app.leave(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: opts.Event,
	})
}

func (app *App) Connect(info ClientInfo, transport Transport) (*Client, error) {
	if transport == nil {
		return nil, IncorrectTransportErr
	}
	if info.ID == NilId {
		info.ID = app.generateClientID()
	}
	client, err := app.pubsub.Connect(info, transport)
	if err == nil {
		app.events.Emit(events.Event{
			Data: client,
			Type: ConnectEvent,
		})
	}
	return client, err
}

func (app *App) DisconnectClient(client *Client) {
	if client == nil {
		return
	}
	app.pubsub.DisconnectClient(client)
	app.events.Emit(events.Event{
		Data: client.id,
		Type: DisconnectEvent,
	})
}

func (app *App) Disconnect(clientId string) {
	app.pubsub.Disconnect(clientId)
	app.events.Emit(events.Event{
		Data: clientId,
		Type: DisconnectEvent,
	})
}

func (app *App) Server(config ServerConfig) Server {
	return &appServer{
		broker:        config.Broker,
		cleanInterval: config.CleanInterval,
		dataHandler:   config.DataHandler,
		starter:       make(chan struct{}, 1),
		App:           app,
	}
}

func NewApp(config AppConfig) *App {
	app := &App{
		id:       config.ID,
		events:   events.NewSource(),
		messages: config.Messages,
	}
	app.pubsub = newPubsub(app, config.PubSub)
	return app
}