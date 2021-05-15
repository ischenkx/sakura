package notify

import (
	"github.com/RomanIschenko/notify/internal/emitter"
	"github.com/RomanIschenko/notify/internal/events"
	"github.com/RomanIschenko/notify/internal/pubsub"
)

type (
	SubscribeClientOptions   = pubsub.SubscribeClientOptions
	SubscribeUserOptions     = pubsub.SubscribeUserOptions
	UnsubscribeClientOptions = pubsub.UnsubscribeClientOptions
	UnsubscribeUserOptions   = pubsub.UnsubscribeUserOptions
	ConnectOptions           = pubsub.ConnectOptions
	DisconnectOptions        = pubsub.DisconnectOptions
	ChangeLog                = *pubsub.ChangeLog
	Priority                 = events.Priority
	EventsCodec              = emitter.EventsCodec
	IDs                      = []string
)
