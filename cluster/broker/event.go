package broker

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Name, Source, ID, BrokerID string
	Time int64
	Data []byte
}


func NewEvent(name string, data []byte) Event {
	return Event{
		Name: name,
		Time: time.Now().UnixNano(),
		ID: uuid.New().String(),
		Data: data,
	}
}