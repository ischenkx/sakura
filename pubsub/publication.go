package pubsub

import (
	"github.com/google/uuid"
	"time"
)

type Publication struct {
	Data []byte
	ID string
	// unix nano
	Time int64
}

func NewPublication(data []byte) Publication {
	return Publication{
		Data: data,
		ID: uuid.New().String(),
		Time: time.Now().UnixNano(),
	}
}