package notifier

import events "github.com/RomanIschenko/notify/event_pubsub"

const (
	BrokerSend events.EventType = "broker_send"
	BrokerJoin					= "broker_join"
	BrokerLeave 				= "broker_leave"
	BrokerAppUp					= "app_up"
	BrokerAppDown				= "app_down"
	BrokerInstanceUp			= "broker_instance_up"
	BrokerInstanceDown			= "broker_instance_down"
)
