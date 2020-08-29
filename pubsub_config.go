package notify

import "time"

type PubSubConfig struct {
	ClientTTL time.Duration
	ClientMessageBufferSize int
}
