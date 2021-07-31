package common

type EventType int

const (
	SystemEvent EventType = iota + 1
	CustomEvent
)

type Event struct {
	Typ  EventType
	From string
	Name string
	Data interface{}
}