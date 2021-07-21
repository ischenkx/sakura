package message

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	id string
	Data []byte
	TimeStamp int64
	Meta interface{}
}

func (m Message) ID() string {
	return m.id
}

func New(data []byte) Message {
	return Message{
		id: uuid.New().String(),
		TimeStamp: time.Now().UnixNano(),
		Data: data,
	}
}
