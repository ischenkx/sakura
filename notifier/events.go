package notifier

import events "github.com/RomanIschenko/notify/event_pubsub"

const (
	BrokerSend events.EventType = "broker_send"
	BrokerJoin					= "broker_join"
	BrokerLeave 				= "broker_leave"
)
