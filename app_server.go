package notify

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify/broker"
	"github.com/RomanIschenko/notify/events"
	"io"
	"io/ioutil"
	"time"
)

type DataHandler func(client *Client, r io.Reader) error

type AppServer struct {
	auth          Auth
	broker        broker.Broker
	cleanInterval time.Duration
	dataHandler	  DataHandler
	starter       chan struct{}
	*App
}

func (app *AppServer) DisconnectClient(client *Client) {
	app.pubsub.DisconnectClient(client)
	app.events.Emit(events.Event{
		Data: client.id, Type: Disconnect,
	})
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
	if app.dataHandler == nil {
		ioutil.ReadAll(r)
		return nil
	}
	return app.dataHandler(client, r)
}

func (app *AppServer) runBroker(ctx context.Context) {
	if app.broker == nil {
		return
	}
	brokerChannel := app.broker.Subscribe()
	defer app.broker.Unsubscribe(brokerChannel)

	appChannel := app.events.Listen()
	defer app.events.Close(appChannel)

	for {
		select {
		case appEvent := <-appChannel:
			switch appEvent.Type {
			case Join, Leave, Send:
				app.broker.Publish(broker.Message{
					Data:     appEvent.Data,
					AppID:    app.id,
					Event:    appEvent.Type,
				})
			}
		case mes := <-brokerChannel:
			if mes.AppID != app.id {
				continue
			}
			switch mes.Event {
			case Send:
				if opts, ok := mes.Data.(SendOptions); ok {
					go app.send(opts)
				}
			case Join:
				if opts, ok := mes.Data.(JoinOptions); ok {
					go app.join(opts)
				}
			case Leave:
				if opts, ok := mes.Data.(LeaveOptions); ok {
					go app.leave(opts)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (app *AppServer) Run(ctx context.Context) {
	app.starter <- struct{}{}
	defer func() {
		<-app.starter
	}()
	app.runBroker(ctx)

	cleaner := time.NewTicker(app.cleanInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-cleaner.C:
			go app.Clean(ctx)
		}
	}
}

func NewAppServer(app *App, config ServerConfig) *AppServer {
	return &AppServer{
		auth:          config.Auth,
		broker:        config.Broker,
		cleanInterval: config.CleanInterval,
		App:           app,
		starter: 	   make(chan struct{}, 1),
	}
}
