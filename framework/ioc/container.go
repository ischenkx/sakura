package ioc

import (
	"errors"
	"github.com/ischenkx/notify/internal/utils"
	"log"
	"reflect"
	"unsafe"
)

const ImportPath = "github.com/ischenkx/notify/framework/ioc"


// Container is not thread safe
type Container struct {
	entries []Entry
	consumers []Consumer
}

// ConsumersRange loops through consumers
func (c *Container) ConsumersRange(f func(Consumer)) {
	if f == nil {return}
	for _, consumer := range c.consumers {
		f(consumer)
	}
}

// Range loop through dependencies
func (c *Container) Range(f func(Entry)) {
	if f == nil {return}
	for _, entry := range c.entries {
		f(entry)
	}
}

func (c *Container) Get(opts ...interface{}) (Entry, bool) {
	var label string
	var typ reflect.Type

	for _, opt := range opts {
		switch val := opt.(type) {
		case LabelOption:
			label = val.data
		case TypeOption:
			typ = val.typ
		}
	}

	if typ == nil && label == "" {
		return Entry{}, false
	}

	for _, entry := range c.entries {
		if (entry.Label == label || label == "") && (matchTypes(typ, entry.typ) || typ == nil) {
			return entry, true
		}
	}

	return Entry{}, false
}

func (c *Container) Inject(data interface{}, opts ...interface{}) error {
	var label string
	for _, opt := range opts {
		switch o := opt.(type) {
		case LabelOption:
			label = o.data
		}
	}
	return c.addEntry(NewEntry(data, label))
}

func (c *Container) addEntry(entry Entry) error {
	if _, ok := c.Get(WithLabel(entry.Label), WithType(entry.typ)); ok {
		return errors.New("already exists")
	}

	if con, ok := c.FindConsumer(entry.typ); ok {
		entry.Value = con
	}

	c.entries = append(c.entries, entry)
	c.initializeConsumers()
	return nil
}

// AddConsumer returns actual value of a consumer
func (c *Container) AddConsumer(con Consumer) (interface{}, error) {
	if con.Object == nil {
		return nil, errors.New("consumer object can not be nil")
	}
	for _, consumer := range c.consumers {
		if consumer.Object == con.Object {
			return nil, errors.New("such consumer already exists")
		}
	}
	if e, ok := c.Get(WithType(reflect.TypeOf(con.Object))); ok {
		con.Object = e.Value
	}
	c.consumers = append(c.consumers, con)
	c.initializeConsumer(&c.consumers[len(c.consumers)-1])
	return con.Object, nil
}

func (c *Container) initializeConsumer(con *Consumer) {
	consumer := con.Object
	fields := con.DependencyMapping
	t := reflect.TypeOf(consumer).Elem()
	val := reflect.ValueOf(consumer).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		lab, ok := fields[field.Name]
		if !ok {
			continue
		}
		e, ok := c.Get(WithType(field.Type), lab)
		if !ok {
			log.Println("failed to find by label and type:", lab, field.Name)
			continue
		}
		reflect.NewAt(field.Type, unsafe.Pointer(val.Field(i).UnsafeAddr())).
			Elem().
			Set(reflect.ValueOf(e.Value))
	}
}

func (c *Container) initializeConsumers() {
	for i := range c.consumers {
		c.initializeConsumer(&c.consumers[i])
	}
}

func (c *Container) FindConsumer(p reflect.Type) (interface{}, bool) {
	for _, consumer := range c.consumers {
		if utils.CompareTypes(p, reflect.TypeOf(consumer)) {
			return consumer, true
		}
	}
	return nil, false
}

func NewEntry(data interface{}, label string) Entry {
	return Entry{
		Label: label,
		Value: data,
		typ:   reflect.TypeOf(data),
	}
}

func New() *Container {
	return &Container{}
}