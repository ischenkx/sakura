package notify

type BrokerEvent int

const (
	Send BrokerEvent = iota + 1
	Join
	Leave
	Connect
	Disconnect
)

type BrokerMessage struct {
	Data interface{}
	AppID string
	Event BrokerEvent
}

type Broker interface {
	Handle(func(message BrokerMessage))
	Send(message BrokerMessage)
}

