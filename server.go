package notify

import (
	"github.com/RomanIschenko/notify/pubsub"
)

type ResolvedConnection struct {
	Client *pubsub.Client
	Err    error
}

type IncomingConnection struct {
	Resolver chan ResolvedConnection
	AuthData string
	Opts     pubsub.ConnectOptions
}

type IncomingData struct {
	Client *pubsub.Client
	Payload []byte
}

type Server interface {
	Accept() <-chan IncomingConnection
	Incoming() <-chan IncomingData
	Inactive() <-chan *pubsub.Client
}
