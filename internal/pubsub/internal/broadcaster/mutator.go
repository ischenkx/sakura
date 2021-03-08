package broadcaster

import (
	"github.com/RomanIschenko/notify/internal/pubsub/internal/mutator"
	"sync"
)

var mutatorsPool = &sync.Pool{
	New: func() interface{} {
		return mutator.New()
	},
}

type Mutator struct {
	mutator.Mutator
	closed	  bool
	timestamp int64
	b         *Broadcaster
}

func (m *Mutator) commit() {
	if m.b != nil {
		m.b.commit(m)
	}
}

func (m *Mutator) Timestamp() int64 {
	return m.timestamp
}

func (m *Mutator) Close() {
	if !m.closed {
		m.commit()
		m.Mutator.Reset()
		mutatorsPool.Put(m.Mutator)
		m.Mutator = nil
		m.b = nil
		m.closed = true
	}
}
