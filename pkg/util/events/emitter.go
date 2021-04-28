package events

import (
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pubsub/message"
	"reflect"
	"sync"
)

// Emitter helps you send and receive different events
// The event encoding algorithm:
// 1) 1st byte is a length of event name
// 2) from the second to 1+<NAME_LENGTH> is the name of event
// 3) other bytes are json encoded data
type Emitter struct {
	app             *notify.App
	appReflectValue reflect.Value
	handlers        map[string]*handler
	codec           Codec
	events          *notify.Events
	rv reflect.Value
	mu sync.RWMutex
}

func (e *Emitter) Emit(ev Event) {
	bts, err := e.codec.Marshal(ev.Data)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	data := make([]byte, 0, 1+len(ev.Name)+len(bts))
	data = append(data, byte(len(ev.Name)))
	data = append(data, []byte(ev.Name)...)
	data = append(data, bts...)

	e.app.Action().
		WithClients(ev.Clients...).
		WithUsers(ev.Users...).
		WithTopics(ev.Topics...).
		WithTimeStamp(ev.TimeStamp).
		WithMetaInfo(ev.MetaInfo).
		Publish(message.New(data))
}

func (e *Emitter) On(name string, handler interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[name] = newHandler(handler, e.codec)
}

func (e *Emitter) processIncomingEvent(a *notify.App, client notify.Client, data []byte) {
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
	hnd.call(e.appReflectValue, e.rv, client, jsonEvent)
}

func (e *Emitter) Close() {
	e.events.Close()
}

// Creates new instance of emitter. Default codec is JSONCodec
func NewEmitter(app *notify.App, codec Codec) *Emitter {
	if codec == nil {
		codec = JSONCodec{}
	}
	emitter := &Emitter{
		app:      app,
		handlers: map[string]*handler{},
		mu:       sync.RWMutex{},
		codec:    codec,
	}
	emitter.appReflectValue = reflect.ValueOf(app)
	emitter.rv = reflect.ValueOf(emitter)
	emitter.events = app.Events(notify.PluginPriority)
	emitter.events.OnMessage(emitter.processIncomingEvent)
	return emitter
}
