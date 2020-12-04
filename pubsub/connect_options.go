package pubsub

import (
	"github.com/RomanIschenko/notify/pubsub/client_id"
	"github.com/RomanIschenko/notify/pubsub/transport"
)

type ConnectOptions struct {
	Transport transport.Transport
	ID        clientid.ID
	Time      int64
}

