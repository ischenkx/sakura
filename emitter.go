package notify

import (
	"errors"
	"github.com/RomanIschenko/notify/internal/utils"
	"reflect"
	"sync"
)

type EventsCodec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type eventHandler struct {
	rawHandler  interface{}
	handlerVal  reflect.Value
	appIndex    int
	clientIndex int
	dataIndex   int
	codec EventsCodec
	dataType    reflect.Type
}

func (h *eventHandler) decodeData(raw []byte) (interface{}, bool) {
	parseDst := reflect.New(h.dataType).Interface()
	if err := h.codec.Unmarshal(raw, parseDst); err == nil {
		return reflect.ValueOf(parseDst).Elem().Interface(), true
	}
	return nil, false
}

func (h *eventHandler) call(app *App, client Client, data interface{}) {
	argsUsed := 0
	var args [3]reflect.Value
	if h.appIndex >= 0 {
		argsUsed++
		args[h.appIndex] = reflect.ValueOf(app)
	}
	if h.clientIndex >= 0 {
		argsUsed++
		args[h.clientIndex] = reflect.ValueOf(client)
	}
	if h.dataIndex >= 0 {
		argsUsed++
		args[h.dataIndex] = reflect.ValueOf(data)
	}
	h.handlerVal.Call(args[:argsUsed])
}

func newHandler(hnd interface{}, codec EventsCodec) *eventHandler {
	t := reflect.TypeOf(hnd)
	if t.NumIn() > 4 {
		panic("eventHandler can't have more than four arguments")
	}

	appType := reflect.TypeOf(&App{})
	clientType := reflect.TypeOf((*Client)(nil)).Elem()
	handlerVal := reflect.ValueOf(hnd)
	h := &eventHandler{
		rawHandler:  hnd,
		handlerVal:  handlerVal,
		appIndex:    -1,
		clientIndex: -1,
		dataIndex:   -1,
		codec: codec,
	}

	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		if utils.CompareTypes(paramType, appType) {
			if h.appIndex >= 0 {
				panic("two apps in one eventHandler")
			}
			h.appIndex = i
		} else if utils.CompareTypes(paramType, clientType) {
			if h.clientIndex >= 0 {
				panic("two clients in one eventHandler")
			}
			h.clientIndex = i
		} else if h.dataIndex < 0 {
			h.dataType = paramType
			h.dataIndex = i
		} else {
			panic("error in eventHandler signature")
		}
	}
	return h
}

type emitter struct {
	app      *App
	handlers map[string]*eventHandler
	codec    EventsCodec
	mu       sync.RWMutex
}

func (e *emitter) encodeData(event string, data interface{}) ([]byte, error) {
	bts, err := e.codec.Marshal(data)

	if err != nil {
		return nil, err
	}
	length := len(event)
	if length >= 256 {
		return nil, errors.New("event name length must be less or equal than 255")
	}
	marshalledData := make([]byte, 1+length+len(bts))
	marshalledData[0] = byte(length)
	copy(marshalledData[1:], event)
	copy(marshalledData[1+length:], bts)
	return marshalledData, nil
}

func (e *emitter) decodeData(data []byte) (string, []byte, error) {
	if len(data) < 2 {
		return "", nil, errors.New("failed to parse event #1")
	}

	length := data[0]
	if len(data) < 1+int(length) {
		return "", nil, errors.New("failed to parse event #2")
	}
	name := data[1 : 1+length]
	event := data[1+length:]
	return string(name), event, nil
}

func (e *emitter) deleteHandler(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.handlers, name)
}

func (e *emitter) handle(name string, hnd interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	e.handlers[name] = newHandler(hnd, e.codec)
}

func (e *emitter) getHandler(name string) (*eventHandler, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	handler, ok := e.handlers[name]
	return handler, ok
}
