package redaptor

import "github.com/ischenkx/swirl"

type client struct {
	adapter *Adapter
	localClient swirl.Client
}

func (c *client) ID() string {
	panic("implement me")
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

func newClient(a *Adapter, lc swirl.Client) swirl.Client {
	return &client{
		adapter:     a,
		localClient: lc,
	}
}
