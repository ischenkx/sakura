package broker

type Subscription interface {
	Close()
	Channel() <-chan Message
}