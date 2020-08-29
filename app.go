package notify

import (
	"context"
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify/events"
	"github.com/google/uuid"
	"log"
	"time"
)

var (
	IncorrectTransportErr = errors.New("transport is nil")
)

type App struct {
	id string
	events *events.Listener
	pubsub *PubSub
	messages MessageStorage
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

func (app *App) Events() *events.Listener {
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
	app.events.Emit(events.Event{
		Data: opts,
		Type: Send,
	})
	if app.messages != nil {
		app.messages.Store(app.id, opts.Message)
	}
}

func (app *App) Join(opts JoinOptions) {
	app.join(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: Join,
	})
}

func (app *App) Leave(opts LeaveOptions) {
	app.leave(opts)
	app.events.Emit(events.Event{
		Data: opts,
		Type: Leave,
	})
}

func (app *App) connect(info ClientInfo, transport Transport) (*Client, error) {
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
			Type: Connect,
		})
	}
	return client, err
}

func (app *App) Disconnect(clientId string) {
	app.pubsub.Disconnect(clientId)
	app.events.Emit(events.Event{
		Data: clientId,
		Type: Disconnect,
	})
}

func NewApp(config AppConfig) *App {
	app := &App{
		id:       config.ID,
		events:   events.NewListener(),
		messages: config.Messages,
	}
	app.pubsub = newPubsub(app, config.PubSub)
	return app
}