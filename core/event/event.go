package event

type Event struct {
	Name string
	Data []byte
}

func New(name string, data []byte) Event {
	return Event{
		Name: name,
		Data: data,
	}
}
