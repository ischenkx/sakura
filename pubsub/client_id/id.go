package clientid

import (
	"errors"
	"github.com/RomanIschenko/notify/internal"
	"github.com/google/uuid"
)

//length of uuid
const Length = 36
//it is placed between ClientClientID and UserClientID
const Separator = "-"

type ID string

func (id ID) Parse() (string, string, error) {
	if len(id) < Length {
		return "", "", errors.New("failed to parse id")
	}

	client := string(id[:Length])
	user := ""

	if len(id) > Length {
		user = string(id[(Length +1):])
	}

	return client, user, nil
}

func (id ID) Hash() (int, error) {
	if len(id) < Length {
		return -1, errors.New("failed to hash")
	}
	bts := []byte(id)
	if len(bts) > Length+1 {
		return internal.Hash(bts[Length:]), nil
	}
	return internal.Hash(bts), nil
}

func New(user, client string) ID {
	sep := ""
	if len(user) > 0 {
		sep = Separator
	}
	if len(client) != Length {
		client = uuid.New().String()
	}
	return ID(client+sep+user)
}

func (id ID) String() string {
	return string(id)
}

func (id ID) User() string {
	if len(id) > Length+1 {
		return string(id[Length+1:])
	}
	return ""
}