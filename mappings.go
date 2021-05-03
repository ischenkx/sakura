package notify

import (
	"github.com/RomanIschenko/notify/internal/events"
	"github.com/RomanIschenko/notify/pubsub"
)

type (
	SubscribeOptions   = pubsub.SubscribeOptions
	UnsubscribeOptions = pubsub.UnsubscribeOptions
	ConnectOptions     = pubsub.ConnectOptions
	DisconnectOptions  = pubsub.DisconnectOptions
	ChangeLog          = pubsub.ChangeLog
	IDs                = []string
	Priority           = events.Priority
)

type Client struct {
	pubsub.Client
}