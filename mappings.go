package notify

import (
	"github.com/ischenkx/notify/internal/emitter"
	"github.com/ischenkx/notify/internal/events"
	"github.com/ischenkx/notify/internal/pubsub"
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
