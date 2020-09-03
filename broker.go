package notify

type Broker interface {
	Subscribe() BrokerSubscription
	Publish(BrokerMessage) error
	ID() string
}