package notify

import (
	"github.com/RomanIschenko/notify/pubsub/client_id"
)

type Auth interface {
	Authorize(string) (clientid.ID, error)
	Register(id clientid.ID) (string, error)
}
