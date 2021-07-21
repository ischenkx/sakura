package redaptor

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ischenkx/swirl"
	"github.com/ischenkx/swirl/pkg/adapter/redis/internal/event"
	"github.com/ischenkx/swirl/pkg/adapter/redis/internal/format"
	"github.com/ischenkx/swirl/pkg/adapter/redis/internal/subscription"
	"github.com/ischenkx/swirl/pkg/adapter/redis/internal/util"
	"time"
)

type worker struct {
	adapter       *Adapter
	subscriptions *subscription.Controller
	pubsub        *redis.PubSub
	ctx           context.Context
}

func (w *worker) processEvents() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case mes := <-w.pubsub.Channel():
			fmt.Println(mes)
		}
	}
}

func (w *worker) bindApp(ctx context.Context) {
	// todo
}

func (w *worker) run(goroutines int) {
	w.bindApp(w.ctx)
	w.subscriptions.Subscribe(w.ctx, "global")
	go w.subscriptions.RunCleaner(w.ctx, int64(time.Minute), int64(time.Minute))
	for i := 0; i < goroutines; i++ {
		go w.processEvents()
	}
}

func newWorker(ctx context.Context, adapter *Adapter) *worker {
	pubsub := adapter.redis.Subscribe(ctx)
	w := &worker{
		adapter:       adapter,
		subscriptions: subscription.NewController(1024, pubsub),
		pubsub:        pubsub,
		ctx:           ctx,
	}
	return w
}

type Adapter struct {
	redis redis.UniversalClient
	app *swirl.App
	uid   string
}

func (a *Adapter) send(ctx context.Context, channel string, ev event.Event) {
	a.redis.Publish(ctx, channel, ev)
}

func (a *Adapter) newEvent(typ uint, name string, data interface{}) event.Event {
	return event.Event{
		Typ:  typ,
		From: a.uid,
		Name: name,
		Data: data,
	}
}

func (a *Adapter) handleSystemEvent(ev event.Event) {
	switch ev.Name {

	}
}

func (a *Adapter) handleCustomEvent(ev event.Event) {

}

func (a *Adapter) handleEvent(ev event.Event) {
	if ev.From == a.uid {
		return
	}

	switch ev.Typ {
	case event.CustomEvent:
		a.handleCustomEvent(ev)
	case event.SystemEvent:
		a.handleSystemEvent(ev)
	}
}

func (a *Adapter) initEvents() {
	events := a.app.Events(swirl.PluginPriority)

	events.OnEvent(func(app *swirl.App, c swirl.Client, options swirl.EventOptions) {
		if s, ok := util.TryJSON(options); ok {
			a.send(context.Background(), format.Client(c.ID()), a.newEvent(event.CustomEvent, "clientEvent", s))
		}
	})

	events.OnConnect(func(app *swirl.App, options swirl.ConnectOptions, c swirl.Client) {
		if s, ok := util.TryJSON(options); ok {
			a.send(context.Background(), format.Client(c.ID()), a.newEvent(event.SystemEvent, "clientConnect", s))
		}
	})

	events.OnReconnect(func(app *swirl.App, options swirl.ConnectOptions, c swirl.Client) {

	})

	events.OnDisconnect(func(app *swirl.App, c swirl.Client) {

	})

	events.OnClientSubscribe(func(app *swirl.App, c swirl.Client, s string, i int64) {

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

func (a *Adapter) Client(id string) swirl.Client {
	return &client{
		adapter:     a,
		localClient: a.app.Client(id, swirl.LocalFlag{}),
	}
}

func (a *Adapter) User(id string) swirl.User {
	return &user{
		adapter:   a,
		localUser: a.app.User(id, swirl.LocalFlag{}),
	}
}

func (a *Adapter) Topic(id string) swirl.Topic {
	return &topic{
		adapter: a,
		topic:   a.app.Topic(id, swirl.LocalFlag{}),
	}
}

func (a *Adapter) Emit(options swirl.EventOptions) {
	a.send(context.Background(), "global", a.newEvent(event.CustomEvent, options.Name, options))
}

func (a *Adapter) Handle(event string, handler interface{}) (swirl.Handle, error) {
	panic("implement me")
}

func (a *Adapter) Metrics() swirl.Metrics {
	panic("implement me")
}

func (a *Adapter) Init(app *swirl.App) {
	a.app = app
}

func (a *Adapter) Start(ctx context.Context) {
	newWorker(ctx, a).run(4)
}

func New(redis redis.UniversalClient) *Adapter {
	adapter := &Adapter{
		redis: redis,
		uid:   uuid.New().String(),
	}
	return adapter
}
