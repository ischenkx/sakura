package pubsub

import (
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify/internal"
	"github.com/google/uuid"
)

// client id is generated with the next algorithm:
// if we have user id and common id then we concatenate them
// using the following pattern "{userId}:{id}"
// else we just use client id

// client id must have constant length to determine whether
// id that we are working with is the concatenation of id and user id
// or just id of a client.

const ClientIDLength = 36

func GetUserID(clientID string) string {
	if len(clientID) < ClientIDLength + 1 {
		return ""
	}
	return clientID[ClientIDLength+1:]
}

func NewClientID(userId string) string {
	id := uuid.New().String()
	if len(userId) == 0 {
		return id
	}
	return fmt.Sprintf("%s:%s", id, userId)
}

func ParseClientID(id string) (clientID string, userID string, err error) {
	if len(id) < ClientIDLength {
		return "", "", errors.New("incorrect client id")
	}
	var userId string
	if len(id) > ClientIDLength {
		userId = id[ClientIDLength+1:]
	}
	return id, userId, nil
}

func HashClientID(id string) (int, error) {
	if len(id) < ClientIDLength {
		return -1, errors.New("id length is incorrect")
	}
	if len(id) > ClientIDLength {
		id = id[ClientIDLength:]
	}
	return internal.Hash([]byte(id)), nil
}