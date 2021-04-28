package common

import (
	"github.com/RomanIschenko/notify/pubsub/protocol"
)

type Flusher interface {
	Flush(protocol protocol.Protocol)
}
