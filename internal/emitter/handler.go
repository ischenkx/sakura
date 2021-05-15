package emitter

import (
	"github.com/RomanIschenko/notify/internal/utils"
	"reflect"
)

// Handler is a function that looks like func(app, client, data)
type Handler struct {
	rawHandler                       interface{}
	handlerVal                       reflect.Value
	dataType                         reflect.Type
	appIndex, clientIndex, dataIndex int
	codec                            EventsCodec
}

func (h *Handler) GetData(raw []byte) (interface{}, bool) {
	if h.dataIndex < 0 {
		return nil, true
	}
	parseDst := reflect.New(h.dataType).Interface()
	if err := h.codec.Unmarshal(raw, parseDst); err == nil {
		return reflect.ValueOf(parseDst).Elem().Interface(), true
	}
	return nil, false
}

func (h *Handler) Call(app interface{}, client interface{}, data interface{}) {
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

func newHandler(hnd interface{}, codec EventsCodec, appType, clientType reflect.Type) *Handler {
	t := reflect.TypeOf(hnd)
	if t.NumIn() > 4 {
		panic("eventHandler can't have more than four arguments")
	}

	handlerVal := reflect.ValueOf(hnd)

	h := &Handler{
		rawHandler:  hnd,
		handlerVal:  handlerVal,
		appIndex:    -1,
		clientIndex: -1,
		dataIndex:   -1,
		codec:       codec,
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
			if utils.IsInterface(paramType) {
				if utils.IsEmptyInterface(paramType) {
					h.dataType = paramType
					h.dataIndex = i
				} else {
					panic("handlers cannot accept non-empty interface data types")
				}
			} else {
				h.dataType = paramType
				h.dataIndex = i
			}
		} else {
			panic("error in eventHandler signature")
		}
	}
	return h
}
