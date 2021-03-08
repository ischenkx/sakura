package common

import "github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/batch"

type Flusher interface {
	Flush(*batch.Batcher)
}
