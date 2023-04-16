package sakura

import (
	"context"
	"log"
	"sakura/event"
	"sakura/subscription"
	"sakura/util"
)

type User struct {
	id     string
	sakura *Sakura
}

func (user User) ID() string {
	return user.id
}

func (user User) Subscribe(ctx context.Context, topic string) error {
	err := user.sakura.subscriptions.Insert(ctx, subscription.Subscription{
		User:  user.id,
		Topic: topic,
	})
	if err != nil {
		return err
	}

	user.PostEvent(ctx, SubscribeEvent, []byte(topic))
	return nil
}

func (user User) Unsubscribe(ctx context.Context, topic string) error {
	set, err := user.sakura.subscriptions.Select(ctx, subscription.Selector{
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

	user.PostEvent(ctx, UnsubscribeEvent, []byte(topic))
	return nil
}

func (user User) UnsubscribeAll(ctx context.Context) error {
	set, err := user.sakura.subscriptions.Select(ctx, subscription.Selector{User: util.MakePtr(user.id)})
	if err != nil {
		return err
	}

	err = set.Erase(ctx)
	if err != nil {
		return err
	}

	user.PostEvent(ctx, UnsubscribeAllEvent, nil)
	return nil
}

func (user User) Subscriptions(ctx context.Context) ([]string, error) {
	set, err := user.sakura.subscriptions.Select(ctx, subscription.Selector{User: util.MakePtr(user.id)})
	if err != nil {
		return nil, err
	}

	var topics []string
	set.Iter(ctx, func(s subscription.Subscription) bool {
		topics = append(topics, s.Topic)
		return true
	})

	return topics, nil
}

func (user User) PostEvent(ctx context.Context, ev string, data []byte) {
	err := user.sakura.Events().Publish(ctx, user.EventChannel(), event.New(ev, data))
	if err != nil {
		log.Println("failed to post an event:", err)
	}
}

func (user User) EventChannel() string {
	return UserChannel(user.ID())
}
