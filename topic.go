package sakura

import (
	"context"
	"sakura/event"
	"sakura/subscription"
	"sakura/util"
)

type Topic struct {
	id     string
	sakura *Sakura
}

func (topic Topic) ID() string {
	return topic.id
}

func (topic Topic) Subscribers(ctx context.Context) ([]string, error) {
	set, err := topic.sakura.subscriptions.Select(ctx, subscription.Selector{
		Topic: util.MakePtr(topic.id),
	})
	if err != nil {
		return nil, err
	}

	var users []string
	set.Iter(ctx, func(sub subscription.Subscription) bool {
		users = append(users, sub.User)
		return true
	})

	return users, nil
}

func (topic Topic) Erase(ctx context.Context) error {
	set, err := topic.sakura.subscriptions.Select(ctx, subscription.Selector{
		Topic: util.MakePtr(topic.id),
	})
	if err != nil {
		return err
	}

	if err := set.Erase(ctx); err != nil {
		return err
	}

	return topic.PostEvent(ctx, TopicErasureEvent, nil)
}

func (topic Topic) Publish(ctx context.Context, data []byte) error {
	return topic.PostEvent(ctx, PublishEvent, data)
}

func (topic Topic) PostEvent(ctx context.Context, ev string, data []byte) error {
	return topic.sakura.Events().Publish(ctx, topic.EventChannel(), event.New(ev, data))
}

func (topic Topic) EventChannel() string {
	return TopicChannel(topic.ID())
}
