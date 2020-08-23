package notify

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"notify/cleaner"
	"sync/atomic"
	"time"
)

const (
	RunningApp int32 = iota + 1
	UpdatingApp
	StoppedApp
)

var (
	StoppedAppErr = errors.New("app is stopped")
	IncorrectTransportErr = errors.New("transport is nil")
)

type AppConfig struct {
	ID string
	Messages MessagesStorage
	Broker Broker
	controlled bool
	PubSub PubSubConfig
}
type App struct {
	id string
	events *EventsHandler
	pubsub *PubSub
	broker Broker
	messages MessagesStorage
	cleaner *cleaner.Cleaner
	controlled bool
	state int32
}

func (app *App) IsRunning() bool {
	return atomic.LoadInt32(&app.state) == RunningApp
}

func (app *App) switchState(from, to int32, updater func()) {
	if atomic.CompareAndSwapInt32(&app.state, from, UpdatingApp) {
		if updater != nil {
			updater()
		}
		atomic.StoreInt32(&app.state, to)
	} else {
		time.Sleep(time.Microsecond * 10)
		app.switchState(from, to, updater)
	}
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
	if !app.IsRunning() {
		return
	}
	app.pubsub.Send(opts)
}

func (app *App) join(opts JoinOptions) {
	if !app.IsRunning() {
		return
	}
	app.pubsub.Join(opts)
}

func (app *App) leave(opts LeaveOptions) {
	if !app.IsRunning() {
		return
	}
	app.pubsub.Leave(opts)
}

func (app *App) loadMessages(ids []string) ([]Message, []string) {
	if app.messages == nil {
		return nil, nil
	}
	return app.messages.Load(app.id, ids)
}

// highly not recommended to call it
// because it's automatically called by cleaner
// in node or in app
func (app *App) Clean() {
	if !app.IsRunning() {
		return
	}
	start := time.Now()
	ids := app.pubsub.Clean()
	elapsed := time.Since(start)
	log.Printf("Cleaning:\n\tDeleted: %v\n\tTime: %v\n", len(ids), elapsed)
	if app.messages != nil {
		app.messages.Clean()
	}
}

func (app *App) Send(mes MessageOptions) {
	if !app.IsRunning() {
		return
	}
	opts := app.identifyMessage(mes)
	app.send(opts)
	app.events.handle(opts)
	if app.messages != nil {
		app.messages.Store(app.id, opts.Message)
	}
}

func (app *App) Join(opts JoinOptions) {
	if !app.IsRunning() {
		return
	}
	app.join(opts)
	app.events.handle(opts)
}

func (app *App) Leave(opts LeaveOptions) {
	if !app.IsRunning() {
		return
	}
	app.leave(opts)
	app.events.handle(opts)
}

func (app *App) Connect(info ClientInfo) (*Client, error) {
	if !app.IsRunning() {
		return nil, StoppedAppErr
	}
	if info.Transport == nil {
		return nil, IncorrectTransportErr
	}
	if info.ID == NilId {
		info.ID = app.generateClientID()
	}
	return app.pubsub.Connect(info)
}

func (app *App) Disconnect(clientId string) {
	fmt.Println("Disconnecting:", clientId)
	app.pubsub.Disconnect(clientId)
}

func (app *App) Run() {
	app.switchState(StoppedApp, RunningApp, func() {
		if !app.controlled {
			if app.cleaner != nil {
				app.cleaner.Run()
			}
			if app.broker != nil {
				app.broker.Handle(func(mes BrokerMessage) {
					if app == nil {
						return
					}
					if app.id != mes.AppID {
						return
					}
					switch mes.Event {
					case Send:
						if opts, ok := mes.Data.(SendOptions); ok {
							app.send(opts)
						}
					case Join:
						if opts, ok := mes.Data.(JoinOptions); ok {
							app.join(opts)
						}
					case Leave:
						if opts, ok := mes.Data.(LeaveOptions); ok {
							app.leave(opts)
						}
					}
				})
			}

			if app.cleaner != nil {
				app.cleaner.Run()
			}
		}

	})
}

func (app *App) Stop() {
	app.switchState(RunningApp, StoppedApp, func() {
		if !app.controlled {
			if app.cleaner != nil {
				app.cleaner.Stop()
			}
		}
	})
}

func (app *App) Server() Server {
	return (*ServableApp)(app)
}

func NewApp(config AppConfig) *App {
	app := &App{
		id: config.ID,
		events: &EventsHandler{},
		broker: config.Broker,
		messages: config.Messages,
		controlled: config.controlled,
		state: StoppedApp,
	}
	if !config.controlled {
		app.cleaner = cleaner.New(app, time.Second * 10)
	}
	app.pubsub = newPubsub(app, config.PubSub)
	return app
}

type ServableApp App

func (sapp *ServableApp) DisconnectClient(client *Client) {
	app := (*App)(sapp)
	if !app.IsRunning() {
		return
	}
	app.pubsub.DisconnectClient(client)
}

func (sapp *ServableApp) Connect(appID string, info ClientInfo) (*Client, error) {
	app := (*App)(sapp)
	if !app.IsRunning() {
		return nil, StoppedAppErr
	}
	if appID != app.id {
		return nil, errors.New("wrong app")
	}
	return app.Connect(info)
}

func (sapp *ServableApp) Handle(client *Client, r io.Reader) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if r == nil {
		return nil
	}
	app := (*App)(sapp)
	err := app.events.handleData(client, r)
	if err != nil {
		app.pubsub.DisconnectClient(client)
	}
	return err
}


