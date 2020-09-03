package notify

import (
	"github.com/RomanIschenko/notify/message"
)

type Send struct {
	Users 	   []string
	Clients    []string
	Channels   []string
	ToBeStored bool
	message.Message
	EventOptions
}
