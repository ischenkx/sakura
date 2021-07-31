package bba

import "github.com/ischenkx/swirl"

type topicEvents struct {
	
}

func (t *topicEvents) OnClientSubscribe(f func(swirl.Client)) {
	panic("implement me")
}

func (t *topicEvents) OnUserSubscribe(f func(swirl.User)) {
	panic("implement me")
}

func (t *topicEvents) OnClientUnsubscribe(f func(swirl.Client)) {
	panic("implement me")
}

func (t *topicEvents) OnUserUnsubscribe(f func(swirl.User)) {
	panic("implement me")
}

func (t *topicEvents) OnEmit(f func(options swirl.EventOptions)) {
	panic("implement me")
}

func (t *topicEvents) Close() {
	panic("implement me")
}

type topic struct {
	adapter *Adapter
	topic swirl.Topic
}

func (t *topic) Emit(s string, i ...interface{}) {
	panic("implement me")
}

func (t *topic) Users() swirl.IDList {
	panic("implement me")
}

func (t *topic) Clients() swirl.IDList {
	panic("implement me")
}

func (t *topic) Events() swirl.TopicEvents {
	panic("implement me")
}



