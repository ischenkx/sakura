package distributor

import "sync"

type sharderPool struct {
	shards int
	pool sync.Pool
}

func (p *sharderPool) Get() *Sharder {
	return p.pool.Get().(*Sharder)
}

func (p *sharderPool) Put(s *Sharder) {
	p.pool.Put(s)
}

func newSharderPool(shards int) *sharderPool {
	return &sharderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return newSharder(shards)
			},
		},
		shards: shards,
	}
}