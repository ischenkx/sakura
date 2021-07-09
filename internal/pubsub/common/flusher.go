package common

import (
	"github.com/ischenkx/swirl/internal/pubsub/protocol"
)

type Flusher interface {
	Flush(protocol protocol.Protocol)
}
