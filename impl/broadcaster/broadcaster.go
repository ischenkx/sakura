package broadcaster

import (
	"context"
	"github.com/samber/lo"
	"log"
	"sakura"
	"sakura/channels"
)

type Broadcaster struct {
	sakura        *sakura.Sakura
	subscriptions *SubscriptionManager
	users         *UserManager
	pubsub        sakura.PubSub
}

func New(sakura *sakura.Sakura) *Broadcaster {
	return &Broadcaster{
		sakura:        sakura,
		subscriptions: newSubscriptionManager(),
		users:         newUserManager(),
		pubsub:        sakura.Broker().PubSub(),
	}
}

func (broadcaster *Broadcaster) Run(ctx context.Context) error {
	return broadcaster.processEvents(ctx)
}

func (broadcaster *Broadcaster) Connect(ctx context.Context, connection User) error {
	user := broadcaster.sakura.User(connection.ID())

	var newTopics []string

	newTopics = append(newTopics, user.Channel())

	subs, err := user.Subscriptions(ctx)
	if err != nil {
		return err
	}

	newTopics = append(
		newTopics,
		lo.Map(subs, func(sub string, _ int) string { return broadcaster.sakura.Topic(sub).Channel() })...,
	)

	if err := broadcaster.pubsub.Subscribe(ctx, newTopics...); err != nil {
		return err
	}
	broadcaster.users.Add(connection)
	for _, sub := range subs {
		broadcaster.subscriptions.Add(sub, connection.ID())
	}
	return nil
}

func (broadcaster *Broadcaster) Disconnect(ctx context.Context, id string) error {
	broadcaster.users.Delete(id)
	broadcaster.subscriptions.RemoveByUser(id)
	return nil
}

func (broadcaster *Broadcaster) processEvents(ctx context.Context) error {
	channel, err := broadcaster.pubsub.Channel(ctx)
	if err != nil {
		return err
	}

	for message := range channel {
		switch message.Data.Name {
		case sakura.SubscribeEvent:
			user := channels.ParseUser(message.Channel)
			topic := string(message.Data.Data)
			broadcaster.subscriptions.Add(user, topic)
		case sakura.UnsubscribeEvent:
			user := channels.ParseUser(message.Channel)
			topic := string(message.Data.Data)
			broadcaster.subscriptions.Remove(user, topic)
		case sakura.UnsubscribeAllEvent:
			user := channels.ParseUser(message.Channel)
			broadcaster.subscriptions.RemoveByUser(user)
		case sakura.TopicErasureEvent:
			topic := channels.ParseTopic(message.Channel)
			broadcaster.subscriptions.RemoveByTopic(topic)
		case sakura.PublishEvent:
			topic := channels.ParseTopic(message.Channel)
			broadcaster.subscriptions.Iter(topic, func(userID string) {
				if user, ok := broadcaster.users.Get(userID); ok {
					if err := user.Send(ctx, message.Data.Data); err != nil {
						log.Println("failed to send data:", err)
					}
				}
			})
		}
	}

	return nil
}
