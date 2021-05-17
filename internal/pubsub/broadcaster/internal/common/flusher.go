package common

import (
	protocol2 "github.com/ischenkx/notify/internal/pubsub/protocol"
)

type Flusher interface {
	Flush(protocol protocol2.Protocol)
}
