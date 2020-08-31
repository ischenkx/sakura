package notifier

type BrokerSubscription interface {
	Close()
	Channel() <-chan BrokerMessage
}