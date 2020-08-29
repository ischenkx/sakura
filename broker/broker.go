package broker

type Broker interface {
	Subscribe() chan Message
	Unsubscribe(chan Message)
	Publish(Message) error
	ID() string
}