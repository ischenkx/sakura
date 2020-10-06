package publication

import (
	"github.com/google/uuid"
	"time"
)

type Publication struct {
	Data []byte
	ID string
	Time int64
}

func New(data []byte) Publication {
	return Publication{
		Data: data,
		ID: uuid.New().String(),
		Time: time.Now().UnixNano(),
	}
}