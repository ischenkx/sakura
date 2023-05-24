package sakura

import (
	"context"
	"sakura/channels"
	"sakura/common/util"
	"sakura/core/event"
	subscription2 "sakura/core/subscription"
	"strings"
)

func ChannelToTopicID(channel string) string {
	return strings.TrimPrefix(channel, "topic/")
}

type AbstractTopic interface {
	ID() string
	Subscribers(ctx context.Context) ([]string, error)
	Drop(ctx context.Context) error
	Publish(ctx context.Context, data []byte) error
	Channel() string
}

type Topic struct {
	id     string
	sakura *Sakura
}

func (topic Topic) ID() string {
	return topic.id
}

func (topic Topic) Subscribers(ctx context.Context) ([]string, error) {
	set, err := topic.sakura.subscriptions.Select(ctx, subscription2.Selector{
		Topic: util.MakePtr(topic.id),
	})
	if err != nil {
		return nil, err
	}

	var users []string
	set.Iter(ctx, func(sub subscription2.Subscription) bool {
		users = append(users, sub.User)
		return true
	})

	return users, nil
}

func (topic Topic) Drop(ctx context.Context) error {
	set, err := topic.sakura.subscriptions.Select(ctx, subscription2.Selector{
		Topic: util.MakePtr(topic.id),
	})
	if err != nil {
		return err
	}
	if err := set.Erase(ctx); err != nil {
		return err
	}
	return topic.pushEvent(ctx, TopicErasureEvent, nil)
}

func (topic Topic) Publish(ctx context.Context, data []byte) error {
	return topic.pushEvent(ctx, PublishEvent, data)
}

func (topic Topic) Channel() string {
	return channels.FromTopic(topic.id)
}

func (topic Topic) pushEvent(ctx context.Context, ev string, data []byte) error {
	return topic.sakura.Broker().Push(ctx, topic.Channel(), event.New(ev, data))
}

type PluginTopic struct {
	base   AbstractTopic
	sakura *Sakura
}

func (topic PluginTopic) ID() string {
	return topic.base.ID()
}

func (topic PluginTopic) Subscribers(ctx context.Context) ([]string, error) {
	return topic.base.Subscribers(ctx)
}

func (topic PluginTopic) Drop(ctx context.Context) error {
	err := topic.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginBeforeTopicDrop); ok {
			if err := p.BeforeTopicDrop(ctx, topic.sakura, topic.ID()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := topic.base.Drop(ctx); err != nil {
		return err
	}

	err = topic.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginAfterTopicDrop); ok {
			p.AfterTopicDrop(ctx, topic.sakura, topic.ID())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (topic PluginTopic) Publish(ctx context.Context, data []byte) error {
	err := topic.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginBeforePublish); ok {
			if err := p.BeforePublish(ctx, topic.sakura, topic.ID(), data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = topic.base.Publish(ctx, data)
	if err != nil {
		return err
	}

	err = topic.sakura.callPlugins(ctx, func(ctx context.Context, plugin Plugin) error {
		if p, ok := plugin.(PluginAfterPublish); ok {
			p.AfterPublish(ctx, topic.sakura, topic.ID(), data)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (topic PluginTopic) Channel() string {
	return topic.base.Channel()
}
