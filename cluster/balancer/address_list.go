package balancer

import "errors"

type AddressList interface {
	Iter(func(int, Instance) bool)
	GetAt(int) (Instance, error)
	Size() int
}

type addrList []Instance

func (t addrList) Size() int {
	return len(t)
}

func (t addrList) GetAt(i int) (Instance, error) {
	if i > len(t) - 1 {
		return Instance{}, errors.New("index is bigger than size of table")
	}
	return t[i], nil
}

func (t addrList) Iter(f func(int, Instance) bool) {
	for i, inst := range t {
		if !f(i, inst) {
			break
		}
	}
}