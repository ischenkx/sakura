package clientid

import (
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify/pubsub/internal/hasher"
	"github.com/google/uuid"
)

// client id is generated with the next algorithm:
// if we have user id and common id then we concatenate them
// using the following pattern "{userId}:{id}"
// else we just use client id

// client id must have constant length to determine whether
// id that we are working with is the concatenation of id and user id
// or just id of a client.

const Length = 36

func User(clientID string) string {
	if len(clientID) < Length+ 1 {
		return ""
	}
	return clientID[Length+1:]
}

func New(userId string) string {
	id := uuid.New().String()
	if len(userId) == 0 {
		return id
	}
	return fmt.Sprintf("%s:%s", id, userId)
}

func Hash(id string) (int, error) {
	if len(id) < Length {
		return -1, errors.New("id length is incorrect")
	}
	if len(id) > Length {
		id = id[Length:]
	}
	return hasher.Hash([]byte(id)), nil
}