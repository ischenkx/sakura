package message

import "github.com/google/uuid"

type Message struct {
	id string
	Data []byte
	Meta interface{}
}

func (m Message) ID() string {
	return m.id
}

func New(data []byte) Message {
	return Message{
		id: uuid.New().String(),
		Data: data,
	}
}
