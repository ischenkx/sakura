package broker

type Broker interface {
	Subscribe() Subscription
	Publish(Message) error
	ID() string
}