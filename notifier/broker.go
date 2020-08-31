package notifier

type Broker interface {
	Subscribe() BrokerSubscription
	Publish(BrokerMessage) error
	ID() string
}