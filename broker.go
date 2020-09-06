package notify

import "io"

type Broker interface {
	Handle(func(BrokerEvent)) io.Closer
	Emit(BrokerEvent)
	ID() string
}