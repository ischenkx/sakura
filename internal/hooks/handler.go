package hooks

import (
	"errors"
	"github.com/ischenkx/swirl/internal/utils"

	"reflect"
)

type handler struct {
	argTypes []reflect.Type
	callable interface{}
}

func (h handler) Call(args []interface{}) ([]reflect.Value, error) {
	if len(h.argTypes) != len(args) {
		return nil, errors.New("invalid amount of passed arguments")
	}

	inputArgs := make([]reflect.Value, len(h.argTypes))

	for i, arg := range args {
		argTyp := h.argTypes[i]
		if !utils.CompareTypes(argTyp, reflect.TypeOf(arg)) && !utils.Implements(reflect.TypeOf(arg), argTyp) {
			return nil, errors.New("type is not matched")
		}
		inputArgs[i] = reflect.ValueOf(arg)
	}
	return reflect.ValueOf(h.callable).Call(inputArgs), nil
}

func newHandler(f interface{}) handler {
	typ := reflect.TypeOf(f)

	argTypes := make([]reflect.Type, typ.NumIn())

	for i := 0; i < typ.NumIn(); i++ {
		argTypes[i] = typ.In(i)
	}

	return handler{
		argTypes: argTypes,
		callable: f,
	}
}
