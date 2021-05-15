package protocol

import (
	message2 "github.com/RomanIschenko/notify/internal/pubsub/message"
)

type Writer interface {
	Write(message2.Batch, interface{})
}

type Protocol interface {
	Encode(writer Writer, buffer message2.Buffer, metadata interface{})
}