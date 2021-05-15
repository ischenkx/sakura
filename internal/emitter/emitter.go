package emitter

import (
	"errors"
	"reflect"
	"sync"
)

type Emitter struct {
	handlers            map[string]*Handler
	codec               EventsCodec
	appType, clientType reflect.Type
	mu                  sync.RWMutex
}

func (e *Emitter) EncodeRawData(event string, data interface{}) ([]byte, error) {
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

func (e *Emitter) DecodeRawData(data []byte) (string, []byte, error) {
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

func (e *Emitter) DeleteHandler(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.handlers, name)
}

func (e *Emitter) GetHandler(name string) (*Handler, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	handler, ok := e.handlers[name]
	return handler, ok
}

func (e *Emitter) Handle(name string, hnd interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	e.handlers[name] = newHandler(hnd, e.codec, e.appType, e.clientType)
}

func New(codec EventsCodec, appType, clientType reflect.Type) *Emitter {
	return &Emitter{
		handlers: map[string]*Handler{},
		codec:      codec,
		appType:    appType,
		clientType: clientType,
	}
}