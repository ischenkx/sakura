package events

import (
	"context"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pubsub"
	"reflect"
	"sync"
)

// Emitter helps you send and receive different events
// The event encoding algorithm:
// 1) 1st byte is a length of event name
// 2) from second to 1+<NAME_LENGTH> is the name of event
// 3) other bytes are json encoded data
type Emitter struct {
	app *notify.App
	handlers map[string]*handler
	codec Codec
	rv reflect.Value
	mu sync.RWMutex
}

func (e *Emitter) Emit(ev Event) {
	bts, err := e.codec.Marshal(ev.Data)
	if err != nil {
		return
	}

	data := make([]byte, 0, 1+len(ev.Name)+len(bts))
	data = append(data, byte(len(ev.Name)))
	data = append(data, []byte(ev.Name)...)
	data = append(data, bts...)

	e.app.Publish(pubsub.PublishOptions{
		Topics:      ev.Topics,
		Clients:     ev.Clients,
		Users:       ev.Users,
		Message:     data,
		NoBuffering: ev.NoBuffering,
		MetaInfo:    ev.MetaInfo,
		Seq:         ev.Seq,
	})
}

func (e *Emitter) On(name string, handler interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[name] = newHandler(handler, e.codec)
}

func (e *Emitter) processIncomingEvent(a *notify.App, client pubsub.Client, data []byte) {
	name, jsonEvent, err := parseIncomingData(data)
	if err != nil {
		logger.Errorln(err)
		return
	}
	e.mu.RLock()
	hnd, ok := e.handlers[name]
	e.mu.RUnlock()
	if !ok {
		return
	}
	hnd.call(reflect.ValueOf(e.app), e.rv, client, jsonEvent)
}

// Creates new instance of emitter. Default codec is JSONCodec
func NewEmitter(ctx context.Context, app *notify.App, codec Codec) *Emitter {
	if codec == nil {
		codec = JSONCodec{}
	}
	emitter := &Emitter{
		app:      app,
		handlers: map[string]*handler{},
		mu:       sync.RWMutex{},
		codec: codec,
	}
	emitter.rv = reflect.ValueOf(emitter)
	app.Events(ctx).OnMessage(emitter.processIncomingEvent)
	return emitter
}


