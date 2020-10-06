package pubsub

import (
	"errors"
	"github.com/google/uuid"
)

//length of uuid
const IDLength = 36
//it is placed between ClientID and UserID
const IDSeparator = "-"

type ClientID string

func (id ClientID) Parse() (string, string, error) {
	if len(id) < IDLength {
		return "", "", errors.New("failed to parse id")
	}

	client := string(id[:IDLength])
	user := ""

	if len(id) > IDLength {
		user = string(id[(IDLength+1):])
	}

	return client, user, nil
}

func (id ClientID) Hash() (int, error) {
	if len(id) < IDLength {
		return -1, errors.New("failed to hash")
	}
	bts := []byte(id)
	if len(bts) > IDLength+1 {
		return hash(bts[IDLength:]), nil
	}
	return hash(bts), nil
}

func NewClientID(user, client string) ClientID {
	sep := ""
	if len(user) > 0 {
		sep = IDSeparator
	}
	if len(client) != IDLength {
		client = uuid.New().String()
	}
	return ClientID(client+sep+user)
}

func (id ClientID) String() string {
	return string(id)
}

func (id ClientID) User() string {
	if len(id) > IDLength+1 {
		return string(id[IDLength+1:])
	}
	return ""
}