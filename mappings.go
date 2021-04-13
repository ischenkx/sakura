package notify

import "github.com/RomanIschenko/notify/internal/pubsub"

type (
	PublishOptions     = pubsub.PublishOptions
	SubscribeOptions   = pubsub.SubscribeOptions
	UnsubscribeOptions = pubsub.UnsubscribeOptions
	ConnectOptions     = pubsub.ConnectOptions
	DisconnectOptions  = pubsub.DisconnectOptions
	ChangeLog          = pubsub.ChangeLog
	Client			   = pubsub.Client
)
