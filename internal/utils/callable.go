package utils

import (
	"errors"
	"reflect"
)

type Callable struct {
	argTypes []reflect.Type
	callable interface{}
}

func (h Callable) Call(args []interface{}) error {
	inputArgs := make([]reflect.Value, len(h.argTypes))

	foundArgs := 0
	for _, arg := range args {
		for i, argTyp := range h.argTypes {
			if CompareTypes(argTyp, reflect.TypeOf(arg)) {
				if inputArgs[i].IsValid() {
					continue
				}
				inputArgs[i] = reflect.ValueOf(arg)
				foundArgs++
				break
			}
		}
	}



	if foundArgs != len(h.argTypes) {
		return errors.New("callable: failed to find expected arguments")
	}
	reflect.ValueOf(h.callable).Call(inputArgs)
	return nil
}

func NewCallable(f interface{}) Callable {
	typ := reflect.TypeOf(f)

	argTypes := make([]reflect.Type, typ.NumIn())

	for i := 0; i < typ.NumIn(); i++ {
		argTypes[i] = typ.In(i)
	}

	return Callable{
		argTypes: argTypes,
		callable: f,
	}
}
