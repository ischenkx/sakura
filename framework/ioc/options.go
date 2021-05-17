package ioc

import "reflect"

type LabelOption struct {
	data string
}

type TypeOption struct {
	typ reflect.Type
}

func WithLabel(s string) LabelOption {
	return LabelOption{data: s}
}


func WithType(p reflect.Type) TypeOption {
	return TypeOption{p}
}

