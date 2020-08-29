package notify

import (
	"github.com/RomanIschenko/notify/events"
)

const (
	Send events.EventType = "send"
	Join 				  = "join"
	Leave 				  = "leave"
	Connect 			  = "connect"
	Disconnect 			  = "disconnect"
	InstanceUp 			  = "instance_up"
	InstanceDown		  = "instance_down"
	AppUp				  = "app_up"
	AppDown				  = "app_down"
)