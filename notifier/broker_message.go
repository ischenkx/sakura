package notifier

import (
	events "github.com/RomanIschenko/notify/event_pubsub"
)

type BrokerMessage struct {
	Data     interface{}
	AppID    string
	BrokerID string
	Event    events.EventType
	Time     int64
}
