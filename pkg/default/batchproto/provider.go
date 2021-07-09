package batchproto

import (
	"github.com/ischenkx/swirl/internal/pubsub/protocol"
	"sync"
)

type Provider struct {
	pool *sync.Pool
}

func (p Provider) New() protocol.Protocol {
	return p.pool.Get().(Protocol)
}

func (p Provider) Put(proto protocol.Protocol) {
	p.pool.Put(proto)
}

func NewProvider(maxBatchSize int) protocol.Provider {
	return Provider{
		pool: &sync.Pool{
			New: func() interface{} {
				return Protocol{
					batcher: newBatch(maxBatchSize),
				}
			},
		},
	}
}
