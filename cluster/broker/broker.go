package broker

// Broker is a data structure that's supposed to implement
// event-based pubsub
// It's lifetime must be limited by the lifetime of a structure you pass a broker in
type Broker interface {
	// handles incoming events
	Handle(func(Event))
	// publishes events to global space
	Publish([]string, Event) error
	// subscribes specified topics
	Subscribe([]string, int64)
	// unsubscribes specified topics
	Unsubscribe([]string, int64)
	ID() string
}