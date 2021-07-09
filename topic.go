package swirl

import (
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"time"
)

type localTopic struct {
	app *App
	id  string
}

func (t localTopic) Events() TopicEvents {
	return t.app.events.forTopic(t.id)
}

func (t localTopic) Emit(name string, options ...interface{}) {
	var args []interface{}
	var metaInfo interface{}
	var timeStamp int64

	for _, opt := range options {
		switch o := opt.(type) {
		case MetaInfo:
			metaInfo = o
		case TimeStamp:
			timeStamp = time.Time(o).UnixNano()
		case Args:
			args = o
		}
	}

	data, err := t.app.emitter.EncodeRawData(name, args)
	if err != nil {
		t.app.events.callError(EncodingError{
			Reason: err,
			EventOptions: EventOptions{
				Name:      name,
				Args:      args,
				TimeStamp: timeStamp,
				MetaInfo:  metaInfo,
			},
		})
		return
	}
	t.app.pubsub.SendToTopic(t.id, message.New(data))

	t.app.events.callEmit(EmitOptions{
		Topics: []string{t.id},
		EventOptions: EventOptions{
			Name:      name,
			Args:      args,
			TimeStamp: timeStamp,
			MetaInfo:  metaInfo,
		},
	})
}

func (t localTopic) Users() IDList {
	return variadicIdList{
		count: func(list IDList) int {
			return t.app.pubsub.Users().CountTopicUsers(t.id)
		},
		array: func(list IDList) []string {
			ids := make([]string, 0, list.Count())
			t.app.pubsub.Users().IterTopicUsers(t.id, func(id string) {
				ids = append(ids, id)
			})
			return ids
		},
	}
}

func (t localTopic) Clients() IDList {
	return variadicIdList{
		count: func(list IDList) int {
			return t.app.pubsub.CountTopicSubscribers(t.id)
		},
		array: func(list IDList) []string {
			return t.app.pubsub.TopicSubscribers(t.id)
		},
	}
}

type Topic interface {
	Emit(string, ...interface{})
	Users() IDList
	Clients() IDList
	Events() TopicEvents
}
