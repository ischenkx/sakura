package notify

import "github.com/RomanIschenko/notify/pubsub"

type Auth interface {
	Authorize(string) (pubsub.ClientID, error)
	Register(id pubsub.ClientID) (string, error)
}
