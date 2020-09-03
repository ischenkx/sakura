package notify

type BrokerSubscription interface {
	Close()
	Channel() <-chan BrokerMessage
}