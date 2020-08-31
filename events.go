package notify

import events "github.com/RomanIschenko/notify/event_pubsub"

const (
	SendEvent         events.EventType = "send"
	JoinEvent                          = "join"
	LeaveEvent                         = "leave"
	ConnectEvent                       = "connect"
	DisconnectEvent                    = "disconnect"
)