package events

type Handler struct {
	closer chan struct{}
	events chan Event
}
