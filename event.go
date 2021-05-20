package notify

type Event struct {
	Name      string
	Data      []interface{}
	Clients   []string
	Users     []string
	Topics    []string
	Timestamp int64
	Meta      interface{}
}

func newEvent(name string, data []interface{}, opts ...interface{}) Event {
	ev := Event{
		Name: name,
		Data: data,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			ev.Meta = o.Data
		case TimeStampOption:
			ev.Timestamp = o.UnixTime
		}
	}
	return ev
}

type EventError struct {
	Payload []byte
	Error	error
}