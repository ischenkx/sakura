package sakura

import (
	"context"
	"log"
	"sakura/channels"
	"sakura/common/util"
	"sakura/core/event"
	subscription2 "sakura/core/subscription"
)

type AbstractUser interface {
	ID() string
	Subscribe(ctx context.Context, topic string) error
	Unsubscribe(ctx context.Context, topic string) error
	Drop(ctx context.Context) error
	Subscriptions(ctx context.Context) ([]string, error)
	Channel() string
}

type User struct {
	id     string
	sakura *Sakura
}

func (user User) ID() string {
	return user.id
}

func (user User) Subscribe(ctx context.Context, topic string) error {
	err := user.sakura.subscriptions.Insert(ctx, subscription2.Subscription{
		User:  user.id,
		Topic: topic,
	})
	if err != nil {
		return err
	}

	user.postEvent(ctx, SubscribeEvent, []byte(topic))
	return nil
}

func (user User) Unsubscribe(ctx context.Context, topic string) error {
	set, err := user.sakura.subscriptions.Select(ctx, subscription2.Selector{
		User:  util.MakePtr(user.id),
		Topic: util.MakePtr(topic),
	})
	if err != nil {
		return err
	}

	err = set.Erase(ctx)
	if err != nil {
		return err
	}

	user.postEvent(ctx, UnsubscribeEvent, []byte(topic))
	return nil
}

func (user User) Drop(ctx context.Context) error {
	set, err := user.sakura.subscriptions.Select(ctx, subscription2.Selector{User: util.MakePtr(user.id)})
	if err != nil {
		return err
	}

	err = set.Erase(ctx)
	if err != nil {
		return err
	}

	user.postEvent(ctx, UnsubscribeAllEvent, nil)
	return nil
}

func (user User) Subscriptions(ctx context.Context) ([]string, error) {
	set, err := user.sakura.subscriptions.Select(ctx, subscription2.Selector{User: util.MakePtr(user.id)})
	if err != nil {
		return nil, err
	}

	var topics []string
	set.Iter(ctx, func(s subscription2.Subscription) bool {
		topics = append(topics, s.Topic)
		return true
	})

	return topics, nil
}

func (user User) Channel() string {
	return channels.FromUser(user.id)
}

func (user User) postEvent(ctx context.Context, ev string, data []byte) {
	err := user.sakura.Broker().Push(ctx, user.Channel(), event.New(ev, data))
	if err != nil {
		log.Println("failed to post an event:", err)
	}
}

type PluginUser struct {
	base   AbstractUser
	sakura *Sakura
}

func (user PluginUser) ID() string {
	return user.base.ID()
}

func (user PluginUser) Subscribe(ctx context.Context, topic string) error {
	err := user.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginBeforeSubscribe); ok {
			if err := p.BeforeSubscribe(ctx, user.sakura, user.ID(), topic); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := user.base.Subscribe(ctx, topic); err != nil {
		return err
	}

	err = user.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginAfterSubscribe); ok {
			p.AfterSubscribe(ctx, user.sakura, user.ID(), topic)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (user PluginUser) Unsubscribe(ctx context.Context, topic string) error {
	err := user.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginBeforeUnsubscribe); ok {
			if err := p.BeforeUnsubscribe(ctx, user.sakura, user.ID(), topic); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := user.base.Unsubscribe(ctx, topic); err != nil {
		return err
	}

	err = user.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginAfterUnsubscribe); ok {
			p.AfterUnsubscribe(ctx, user.sakura, user.ID(), topic)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (user PluginUser) Drop(ctx context.Context) error {
	err := user.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginBeforeUserDrop); ok {
			if err := p.BeforeUserDrop(ctx, user.sakura, user.ID()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := user.base.Drop(ctx); err != nil {
		return err
	}

	err = user.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginAfterUserDrop); ok {
			p.AfterUserDrop(ctx, user.sakura, user.ID())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (user PluginUser) Subscriptions(ctx context.Context) ([]string, error) {
	return user.base.Subscriptions(ctx)
}

func (user PluginUser) Channel() string {
	return user.base.Channel()
}
