package authmock

import (
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/google/uuid"
)

type Auth struct {}

func (a Auth) Authorize(string) (pubsub.ClientID, error) {
	return pubsub.NewClientID(uuid.New().String(), uuid.New().String()), nil
}

func (a Auth) Register(id pubsub.ClientID) (string, error) {
	return "", nil
}

func New() Auth {
	return Auth{}
}
