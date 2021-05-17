package redis

import (
	"context"
	"encoding/json"
	"github.com/ischenkx/notify"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/common/log"
)

type metaLabel struct {
	labelFromRedisBroker int64
}

func checkMeta(m interface{}) bool {
	_, ok := m.(metaLabel)
	return ok
}

type Broker struct {
	client redis.UniversalClient
}

type event struct {
	Name string
	Data interface{}
}

func (b *Broker) Start(ctx context.Context, app *notify.App) {
	hooks := app.Hooks(notify.PluginPriority)


	pubsub := b.client.Subscribe(ctx)

	format := formatter{}

	// unsubscribe or subscribe deleted topics, users, clients
	hooks.OnChange(func(app *notify.App, log notify.ChangeLog) {
		if ids, ok := format.Clients(log.ClientsDeleted()); ok {
			pubsub.Unsubscribe(ctx, ids...)
		}

		if ids, ok := format.Clients(log.ClientsCreated()); ok {
			pubsub.Subscribe(ctx, ids...)
		}

		if ids, ok := format.Users(log.UsersDeleted()); ok {
			pubsub.Unsubscribe(ctx, ids...)
		}

		if ids, ok := format.Users(log.UsersCreated()); ok {
			pubsub.Subscribe(ctx, ids...)
		}

		if ids, ok := format.Topics(log.TopicsDeleted()); ok {
			pubsub.Unsubscribe(ctx, ids...)
		}

		if ids, ok := format.Topics(log.TopicsCreated()); ok {
			pubsub.Subscribe(ctx, ids...)
		}
	})

	// broadcast emit
	hooks.OnEmit(func(app *notify.App, e notify.Event) {
		if checkMeta(e.Meta) {
			return
		}

		abstractEvent := event{
			Name: e.Name,
			Data: e.Data,
		}

		marshalledEvent, err := json.Marshal(abstractEvent)

		if err != nil {
			log.Info("failed to marshal event:", err)
			return
		}

		pp := b.client.Pipeline()

		for _, c := range e.Clients {
			pp.Publish(ctx, format.Client(c), marshalledEvent)
		}

		for _, u := range e.Users {

		}

		for _, t := range e.Topics {

		}

	})
}

func New() *Broker {
	return &Broker{}
}
