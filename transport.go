package notify

import "github.com/RomanIschenko/notify/message"

type Transport interface {
	Send(message message.Message)
	Close()
}
