package notify

import (
	"github.com/RomanIschenko/notify/pubsub"
	"io"
)

type Broker interface {
	HandlePubsubResult(appID string, result pubsub.Result)
	Handle(func(BrokerEvent)) io.Closer
	Emit(BrokerEvent)
}