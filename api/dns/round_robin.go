package dns

import (
	"errors"
	"fmt"
	"sync/atomic"
)

type roundRobin struct {
	ids []string
	current int32
	table map[string]struct{}
}

func (r *roundRobin) add(id string) {
	if _, ok := r.table[id]; ok {
		return
	}
	r.table[id] = struct{}{}
	r.ids = append(r.ids, id)
}

func (r *roundRobin) del(id string) {
	if _, ok := r.table[id]; !ok {
		fmt.Println("nothing to delete")
		return
	}
	delete(r.table, id)
	l := len(r.ids)
	for i := 0; i < len(r.ids); i++ {
		if r.ids[i] == id {
			r.ids[len(r.ids)-1], r.ids[i] = id, r.ids[len(r.ids)-1]
			r.ids = r.ids[:len(r.ids)-1]
		}
	}

	fmt.Println("might have deleted smth:", l, len(r.ids))
}

func (r *roundRobin) next() (string, error) {
	if len(r.ids) == 0 {
		return "", errors.New("no ids provided")
	}
	i := int(atomic.AddInt32(&r.current, 1))
	if i >= len(r.ids) {
		atomic.CompareAndSwapInt32(&r.current, int32(i), -1)
		i = 0
	}
	s := r.ids[i]
	return s, nil
}

func newRoundRobin() *roundRobin {
	return &roundRobin{
		ids:     []string{},
		table: map[string]struct{}{},
	}
}
