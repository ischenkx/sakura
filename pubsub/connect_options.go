package pubsub

import (
	"github.com/RomanIschenko/notify/pubsub/transport"
)

type ConnectOptions struct {
	Transport transport.Transport
	ID    	  string
	Time      int64
}

