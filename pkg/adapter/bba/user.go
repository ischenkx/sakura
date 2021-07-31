package bba

import "github.com/ischenkx/swirl"

type user struct {
	adapter *adapter
	localUser swirl.User
}

func (u *user) ID() string {
	panic("implement me")
}

func (u *user) Active() bool {
	panic("implement me")
}

func (u *user) Emit(name string, data ...interface{}) {
	panic("implement me")
}

func (u *user) Subscribe(options swirl.SubscribeOptions) {
	panic("implement me")
}

func (u *user) Unsubscribe(options swirl.UnsubscribeOptions) {
	panic("implement me")
}

func (u *user) Disconnect(options swirl.DisconnectOptions) {
	panic("implement me")
}

func (u *user) Clients() []string {
	panic("implement me")
}

func (u *user) Events() swirl.UserEvents {
	panic("implement me")
}

func (u *user) Subscriptions() swirl.SubscriptionList {
	panic("implement me")
}


