package notify

import (
	"github.com/RomanIschenko/notify/message"
)

type SendOptions struct {
	Users 	   []string
	Clients    []string
	Channels   []string
	ToBeStored bool
	message.Message
	Event string
}
