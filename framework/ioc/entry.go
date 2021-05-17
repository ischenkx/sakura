package ioc

import "reflect"

type Entry struct {
	Label string
	Value interface{}
	typ reflect.Type
}
