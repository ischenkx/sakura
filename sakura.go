package sakura

import (
	"context"
	"sakura/core/broker"
	"sakura/core/event"
	"sakura/core/subscription"
)

type Builder struct {
	Subscriptions subscription.Storage
	Broker        broker.Broker[event.Event]
}

func (builder Builder) Build() *Sakura {
	return &Sakura{
		subscriptions: builder.Subscriptions,
		broker:        builder.Broker,
	}
}

type Sakura struct {
	subscriptions subscription.Storage
	broker        broker.Broker[event.Event]
	plugins       []Plugin
}

func (sakura *Sakura) Use(ctx context.Context, plugin Plugin) error {
	if initializer, ok := plugin.(PluginInitializer); ok {
		if err := initializer.Initialize(ctx, sakura); err != nil {
			return err
		}
	}
	sakura.plugins = append(sakura.plugins, plugin)
	return nil
}

func (sakura *Sakura) User(id string) AbstractUser {
	return PluginUser{
		base: User{
			id:     id,
			sakura: sakura,
		},
		sakura: sakura,
	}
}

func (sakura *Sakura) Topic(id string) AbstractTopic {
	return PluginTopic{
		base: Topic{
			id:     id,
			sakura: sakura,
		},
		sakura: sakura,
	}
}

func (sakura *Sakura) Broker() broker.Broker[event.Event] {
	return sakura.broker
}

func (sakura *Sakura) callPlugins(ctx context.Context, caller func(ctx context.Context, plugin Plugin) error) error {
	if caller == nil {
		return nil
	}
	for _, plugin := range sakura.plugins {
		if err := caller(ctx, plugin); err != nil {
			return err
		}
	}
	return nil
}
