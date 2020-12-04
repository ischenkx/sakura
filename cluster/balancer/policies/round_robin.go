package policies

import (
	"github.com/RomanIschenko/notify/cluster/balancer"
	"sync/atomic"
)

type RoundRobin struct {
	index int32
}

func (r *RoundRobin) Next(list balancer.AddressList) (balancer.Instance, error) {
	current := atomic.AddInt32(&r.index, 1)
	if int(current) >= list.Size() {
		atomic.CompareAndSwapInt32(&r.index, current, 0)
		current = 0
	}
	return list.GetAt(int(current))
}