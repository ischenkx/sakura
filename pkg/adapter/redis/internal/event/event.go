package event

import "encoding/json"

const (
	SystemEvent uint = iota + 1
	CustomEvent
)

type Event struct {
	Typ  uint
	From string
	Name string
	Data interface{}
}

func (e Event) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

