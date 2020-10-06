package pubsub

import (
	"github.com/RomanIschenko/notify/pubsub/publication"
)

type SubscribeOptions struct {
	Topics []string
	Clients []string
	Users []string
	Time 	int64
}

type UnsubscribeOptions struct {
	Topics []string
	Clients []string
	Users []string
	All bool
	Time 	int64
}

type PublishOptions struct {
	Topics []string
	Clients []string
	Users []string
	Payload publication.Publication
	Time 	int64
}

type ConnectOptions struct {
	Transport Transport
	ID		  ClientID
	Time 	  int64
}

type DisconnectOptions struct {
	Clients []string
	Users	[]string
	All		bool
	Time 	int64
}