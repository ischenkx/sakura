package protocol

import "github.com/RomanIschenko/notify/pubsub/message"

type Writer interface {
	Write(message.Batch, interface{})
}

type Protocol interface {
	Encode(writer Writer, buffer message.Buffer, metadata interface{})
}