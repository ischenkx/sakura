package emitter

import (
	"errors"
	"fmt"
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

// Handler is a function that looks like func([app,] [client,] args...)
type Handler struct {
	rawHandler            interface{}
	handlerVal            reflect.Value
	appIndex, clientIndex int
	dataTypes			  []reflect.Type
	codec                 EventsCodec
}

func (h *Handler) getZeroArgs() []interface{} {
	args := make([]interface{}, len(h.dataTypes))
	for i, argType := range h.dataTypes {
		args[i] = reflect.New(argType).Interface()
	}
	return args
}

func (h *Handler) DecodeData(data []byte) ([]interface{}, error) {
	args := h.getZeroArgs()
	return args, h.codec.Unmarshal(data, args)
}

func (h *Handler) Call(app interface{}, client interface{}, data ...interface{}) error {

	if len(h.dataTypes) != len(data) {
		return fmt.Errorf("expected %d arguments, but received %d", len(h.dataTypes), len(data))
	}

	length := 0

	if h.appIndex >= 0 {
		length++
	}
	if h.clientIndex >= 0 {
		length++
	}
	length += len(h.dataTypes)

	args := make([]reflect.Value, length)
	argsUsed := 0

	if h.appIndex >= 0 {
		args[h.appIndex] = reflect.ValueOf(app)
		argsUsed++
	}
	if h.clientIndex >= 0 {
		args[h.clientIndex] = reflect.ValueOf(client)
		argsUsed++
	}
	for i := 0; i < len(data); i++ {
		arg := data[i]
		argType := reflect.TypeOf(arg).Elem()
		val := reflect.ValueOf(arg).Elem()
		if !utils.CompareTypes(h.dataTypes[i], argType) {
			return fmt.Errorf("unexpected arg at position %d: exepected %s, but received %s", i+argsUsed+1, h.dataTypes[i].String(), argType.String())
		}
		args[argsUsed+i] = val
	}
	h.handlerVal.Call(args)
	return nil
}

func newHandler(hnd interface{}, codec EventsCodec, appType, clientType reflect.Type) (*Handler, error) {
	t := reflect.TypeOf(hnd)
	if t.NumIn() > 4 {
		return nil, errors.New("eventHandler can't have more than four arguments")
	}

	handlerVal := reflect.ValueOf(hnd)

	h := &Handler{
		rawHandler:  hnd,
		handlerVal:  handlerVal,
		appIndex:    -1,
		clientIndex: -1,
		codec:       codec,
	}
	dataTypesFlag := false

	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		if utils.CompareTypes(paramType, appType) {
			if h.appIndex >= 0 {
				return nil, errors.New("two apps in one eventHandler")
			}

			if dataTypesFlag {
				return nil, errors.New("app parameter can not be placed after data parameters")
			}

			h.appIndex = i
			continue
		} else if utils.CompareTypes(paramType, clientType) {
			if h.appIndex >= 0 {
				return nil, errors.New("two clients in one eventHandler")
			}

			if dataTypesFlag {
				return nil, errors.New("client parameter can not be placed after data parameters")
			}
			h.clientIndex = i
			continue
		}
		dataTypesFlag = true
		h.dataTypes = append(h.dataTypes, paramType)
	}
	return h, nil
}
