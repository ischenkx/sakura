package batchproto

import (
	protocol2 "github.com/RomanIschenko/notify/internal/pubsub/protocol"
	"sync"
)

type Provider struct {
	pool *sync.Pool
}

func (p Provider) New() protocol2.Protocol {
	return p.pool.Get().(Protocol)
}

func (p Provider) Put(proto protocol2.Protocol) {
	p.pool.Put(proto)
}

func NewProvider(maxBatchSize int) protocol2.Provider {
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