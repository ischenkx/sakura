package bba

import (
	"context"
	"fmt"
	"github.com/ischenkx/swirl"
	"github.com/ischenkx/swirl/pkg/adapter/bba/common"
	"github.com/ischenkx/swirl/pkg/adapter/bba/internal/format"
	"github.com/ischenkx/swirl/pkg/adapter/bba/internal/subscription"
	"time"
)

type adapter struct {
	app   *swirl.App
	uid    string
	pubsub common.PubSub
}

type worker struct {
	adapter       *adapter
	app           *swirl.App
	subscriptions *subscription.Controller
	ctx           context.Context
}

func (w *worker) processEvents() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case mes := <-w.adapter.pubsub.Channel():
			fmt.Println(mes)
		}
	}
}

func (w *worker) run(goroutines int) {
	w.subscriptions.Subscribe(w.ctx, "global")
	go w.subscriptions.RunCleaner(w.ctx, int64(time.Minute), int64(time.Minute))
	for i := 0; i < goroutines; i++ {
		go w.processEvents()
	}
}

const (
	clientRecvEvent = "client-recv"
	userClientRecvEvent = "user-client-recv"
	userClientConnectEvent = "user-client-connect"
	clientReconnectEvent = "client-reconnect"
	userClientReconnectEvent = "user-client-reconnect"
	clientDisconnectEvent = "client-disconnect"
	userClientDisconnectEvent = "client-disconnect"
)

func (w *worker) initialize() {
	events := w.app.Events(swirl.PluginPriority)
	go func(ctx context.Context, handle swirl.Handle) {
		select {
		case <-ctx.Done():
			handle.Close()
		}
	}(w.ctx, events)
	background := context.Background()
	events.OnEvent(func(app *swirl.App, c swirl.Client, options swirl.EventOptions) {
		// we received an event from a client, so we should broadcast info about it
		// firstly let's send data to the client personal channel
		clientChannel := format.Client(c.ID())

		w.adapter.send(background,
			clientChannel,
			w.adapter.newEvent(common.SystemEvent, clientRecvEvent, options))

		// ok, now we have to notify a user that's bound to the client
		if c.User().ID() != "" {
			userChannel := format.User(c.User().ID())
			w.adapter.send(
				background,
				userChannel,
				w.adapter.newEvent(common.SystemEvent, userClientRecvEvent, options))
		}
	})

	events.OnConnect(func(app *swirl.App, options swirl.ConnectOptions, c swirl.Client) {
		if c.User().ID() != "" {
			userChannel := format.User(c.User().ID())
			w.adapter.send(
				background,
				userChannel,
				w.adapter.newEvent(common.SystemEvent, userClientConnectEvent, options))
		}
	})

	events.OnReconnect(func(app *swirl.App, options swirl.ConnectOptions, c swirl.Client) {
		clientChannel := format.Client(c.ID())
		payload := options
		w.adapter.send(background,
			clientChannel,
			w.adapter.newEvent(common.SystemEvent, clientReconnectEvent, payload))
		// ok, now we have to notify a user that's bound to the client
		if c.User().ID() != "" {
			userChannel := format.User(c.User().ID())
			w.adapter.send(
				background,
				userChannel,
				w.adapter.newEvent(common.SystemEvent, userClientReconnectEvent, payload))
		}
	})

	events.OnDisconnect(func(app *swirl.App, c swirl.Client) {
		clientChannel := format.Client(c.ID())
		payload := c.ID()
		w.adapter.send(background,
			clientChannel,
			w.adapter.newEvent(common.SystemEvent, clientDisconnectEvent, payload))
		// ok, now we have to notify a user that's bound to the client
		if c.User().ID() != "" {
			userChannel := format.User(c.User().ID())
			w.adapter.send(
				background,
				userChannel,
				w.adapter.newEvent(common.SystemEvent, userClientDisconnectEvent, payload))
		}
	})

	events.OnClientSubscribe(func(app *swirl.App, c swirl.Client, topicID string, ts int64) {
		clientChannel := format.Client(c.ID())
		payload := s
		w.adapter.send(background,
			clientChannel,
			w.adapter.newEvent(common.SystemEvent, clientReconnectEvent, payload))
		// ok, now we have to notify a user that's bound to the client
		if c.User().ID() != "" {
			userChannel := format.User(c.User().ID())
			w.adapter.send(
				background,
				userChannel,
				w.adapter.newEvent(common.SystemEvent, userClientReconnectEvent, payload))
		}
	})

	events.OnClientUnsubscribe(func(app *swirl.App, c swirl.Client, s string, i int64) {

	})

	events.OnUserSubscribe(func(app *swirl.App, u swirl.User, s string, i int64) {

	})

	events.OnUserUnsubscribe(func(app *swirl.App, u swirl.User, s string, i int64) {

	})

	events.OnEmit(func(app *swirl.App, options swirl.EmitOptions) {

	})

	events.OnInactivate(func(app *swirl.App, c swirl.Client) {

	})
}

func (a *adapter) send(ctx context.Context, channel string, ev common.Event) {
	a.pubsub.Publish(ctx, channel, ev)
}

func (a *adapter) newEvent(typ common.EventType, name string, data interface{}) common.Event {
	return common.Event{
		Typ:  typ,
		From: a.uid,
		Name: name,
		Data: data,
	}
}

func (a *adapter) handleSystemEvent(ev common.Event) {
	switch ev.Name {

	}
}

func (a *adapter) handleCustomEvent(ev common.Event) {

}

func (a *adapter) handleEvent(ev common.Event) {
	if ev.From == a.uid {
		return
	}
	switch ev.Typ {
	case common.CustomEvent:
		a.handleCustomEvent(ev)
	case common.SystemEvent:
		a.handleSystemEvent(ev)
	}
}

func (a *adapter) wrap(ctx context.Context) *worker {
	w := &worker{
		adapter:       a,
		app:           a.app,
		subscriptions: subscription.NewController(1024, a.pubsub),
		ctx:           ctx,
	}
	w.initialize()
	return w
}

func (a *adapter) Client(s string) swirl.Client {
	return &client{
		adapter:     a,
		localClient: a.app.Client(s, swirl.LocalFlag{}),
	}
}

func (a *adapter) User(s string) swirl.User {
	return &user{
		adapter:   a,
		localUser: a.app.User(s, swirl.LocalFlag{}),
	}
}

func (a *adapter) Topic(s string) swirl.Topic {
	return &topic{
		adapter: a,
		topic:   a.app.Topic(s, swirl.LocalFlag{}),
	}
}

func (a *adapter) Emit(options swirl.EventOptions) {
	panic("(emit) custom events are not implemented yet")
}

func (a *adapter) Metrics() swirl.Metrics {
	panic("metrics events are not implemented yet")
}

func (a *adapter) Handle(s string, i interface{}) (swirl.Handle, error) {
	panic("custom events are not implemented yet")
}

func (a *adapter) Start(ctx context.Context) {
	a.wrap(ctx).run(4)
}

