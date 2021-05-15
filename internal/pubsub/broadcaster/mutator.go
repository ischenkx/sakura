package broadcaster

import (
	mutator2 "github.com/RomanIschenko/notify/internal/pubsub/mutator"
	"sync"
)

var mutatorsPool = &sync.Pool{
	New: func() interface{} {
		return mutator2.New()
	},
}

type Mutator struct {
	mutator2.Mutator
	closed    bool
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
