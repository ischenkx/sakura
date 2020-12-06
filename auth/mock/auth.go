package authmock

import (
	clientid "github.com/RomanIschenko/notify/pubsub/clientid"
	"github.com/google/uuid"
)

type Auth struct {}

func (a Auth) Authorize(string) (string, error) {
	return clientid.New(uuid.New().String()), nil
}

func (a Auth) Register(id string) (string, error) {
	return "", nil
}

func New() Auth {
	return Auth{}
}
