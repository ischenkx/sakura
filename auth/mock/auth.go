package authmock

import (
	clientid "github.com/RomanIschenko/notify/pubsub/client_id"
	"github.com/google/uuid"
)

type Auth struct {}

func (a Auth) Authorize(string) (clientid.ID, error) {
	return clientid.New(uuid.New().String(), uuid.New().String()), nil
}

func (a Auth) Register(id clientid.ID) (string, error) {
	return "", nil
}

func New() Auth {
	return Auth{}
}
