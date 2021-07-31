package bba

import "github.com/ischenkx/swirl"

type client struct {
	adapter *adapter
	localClient swirl.Client
}

func (c *client) ID() string {
	return c.localClient.ID()
}

func (c *client) Emit(s string, i ...interface{}) {
	panic("implement me")
}

func (c *client) Subscribe(options swirl.SubscribeOptions) {
	panic("implement me")
}

func (c *client) Unsubscribe(options swirl.UnsubscribeOptions) {
	panic("implement me")
}

func (c *client) Disconnect(options swirl.DisconnectOptions) {
	panic("implement me")
}

func (c *client) User() swirl.User {
	panic("implement me")
}

func (c *client) Active() bool {
	panic("implement me")
}

func (c *client) Events() swirl.ClientEvents {
	panic("implement me")
}

func (c *client) Subscriptions() swirl.SubscriptionList {
	panic("implement me")
}
