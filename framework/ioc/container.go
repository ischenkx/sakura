package ioc

import (
	"errors"
	"github.com/RomanIschenko/notify/internal/utils"
	"log"
	"reflect"
	"unsafe"
)

const ImportPath = "github.com/RomanIschenko/notify/framework/ioc"

type FieldMapping map[string]string

type Entry struct {
	Label string
	Value interface{}
	typ reflect.Type
}

type Container struct {
	entries []Entry
	consumers []interface{}
}

func (c *Container) findByType(t reflect.Type, label string) (Entry, bool) {
	for _, entry := range c.entries {
		if utils.CompareTypes(t, entry.typ) && (label == "" || entry.Label == label) {
			return entry, true
		}

		interfaceRelation := false

		if entry.typ.Kind() == reflect.Interface {
			interfaceRelation = interfaceRelation||t.Implements(entry.typ)
		}

		if t.Kind() == reflect.Interface {
			interfaceRelation = interfaceRelation||entry.typ.Implements(t)
		}

		if entry.Label != "" && entry.Label == label && interfaceRelation {
			return entry, true
		}
	}

	return Entry{}, false
}

func (c *Container) Add(entry Entry) error {
	if _, ok := c.findByType(entry.typ, entry.Label); ok {
		return errors.New("already exists")
	}

	c.entries = append(c.entries, entry)

	return nil
}

func (c *Container) tryFindProvidedConsumerDouble(con interface{}) interface{} {
	for _, entry := range c.entries {
		if utils.CompareTypes(reflect.TypeOf(con), entry.typ) {
			return entry.Value
		}
	}
	return con
}

// InitConsumer accepts a pointer to a struct
func (c *Container) InitConsumer(consumer interface{}, fields FieldMapping) {

	consumer = c.tryFindProvidedConsumerDouble(consumer)

	t := reflect.TypeOf(consumer).Elem()
	val := reflect.ValueOf(consumer).Elem()

	c.consumers = append(c.consumers, consumer)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		lab, ok := fields[field.Name]
		if !ok {
			continue
		}
		e, ok := c.findByType(field.Type, lab)
		if !ok {
			log.Println("failed to find by label and type:", lab, field.Name)
			continue
		}

		reflect.NewAt(field.Type, unsafe.Pointer(val.Field(i).UnsafeAddr())).
			Elem().
			Set(reflect.ValueOf(e.Value))
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