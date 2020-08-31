package notifier

import events "github.com/RomanIschenko/notify/events_pubsub"

const (
	BrokerSend events.EventType = "broker_send"
	BrokerJoin					= "broker_join"
	BrokerLeave 				= "broker_leave"
)
