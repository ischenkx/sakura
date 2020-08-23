package notify

type Transport interface {
	Send(message Message)
	Close()
}
