package broker

import (
	"github.com/RomanIschenko/notify/events"
	"time"
)

type Message struct {
	Data     interface{}
	AppID    string
	BrokerID string
	Event    events.EventType
	Time     time.Time
}
