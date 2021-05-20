package notify

import (
	"errors"
	pubsub2 "github.com/ischenkx/notify/internal/pubsub"
)

type Client struct {
	id        string
	app       *App
	rawClient pubsub2.Client
}

func (c *Client) initRawClient() error {
	if c.rawClient == nil {
		rawClient, err := c.app.pubsub.Client(c.id)
		if err != nil {
			return err
		}
		c.rawClient = rawClient
	}

	if !c.rawClient.Valid() {
		return errors.New("underlying client is invalid")
	}

	return nil
}

func (c Client) Emit(name string, data ...interface{}) {
	ev := newEvent(name, data)
	ev.Clients = []string{c.id}
	c.app.Emit(ev)
}

func (c Client) Subscriptions() ([]string, error) {
	return c.app.pubsub.ClientSubscriptions(c.ID())
}

func (c Client) Subscribed(t string) (bool, error) {
	return c.app.pubsub.ClientSubscribed(c.id, t)
}

func (c Client) Subscribe(topic string, opts ...interface{}) error {
	
	subOpts := SubscribeClientOptions{
		ID:        c.ID(),
		Topic:     topic,
	}
	
	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			subOpts.Meta = o.Data
		case TimeStampOption:
			subOpts.TimeStamp = o.UnixTime
		}
	}

	return c.app.SubscribeClient(subOpts)
}

func (c Client) Unsubscribe(topic string, opts ...interface{}) error {
	unsubOpts := UnsubscribeClientOptions{
		ID:        c.ID(),
		Topic:     topic,
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			unsubOpts.Meta = o.Data
		case TimeStampOption:
			unsubOpts.TimeStamp = o.UnixTime
		}
	}

	return c.app.UnsubscribeClient(unsubOpts)
}

func (c Client) UnsubscribeAll(opts ...interface{}) error {
	unsubOpts := UnsubscribeClientOptions{
		ID: c.ID(),
		All: true,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			unsubOpts.Meta = o.Data
		case TimeStampOption:
			unsubOpts.TimeStamp = o.UnixTime
		}
	}
	return c.app.UnsubscribeClient(unsubOpts)
}

func (c Client) ID() string {
	return c.id
}

func (c Client) User() (User, error) {
	if err := c.initRawClient(); err != nil {
		return User{}, err
	}

	userID := c.rawClient.User()

	if userID == "" {
		return User{}, errors.New("no user is bound to the client")
	}

	return User{
		id:  c.rawClient.User(),
		app: nil,
	}, nil
}