package dns

import (
	"sync"
)

type RoundRobin struct {
	rb *roundRobin
	mu sync.Mutex
}

func (r *RoundRobin) Add(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rb.add(id)
}

func (r *RoundRobin) Del(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rb.del(id)
}

func (r *RoundRobin) Next() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rb.next()
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		rb: newRoundRobin(),
	}
}