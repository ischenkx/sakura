package authmock

import (
	"github.com/google/uuid"
)

type Auth struct {}

func (a Auth) Authorize(string) (string, string, error) {
	return uuid.New().String(), uuid.New().String(), nil
}

func (a Auth) Register(id, user string) (string, error) {
	return "", nil
}

func New() Auth {
	return Auth{}
}
