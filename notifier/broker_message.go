package notifier

type BrokerMessage struct {
	Data     interface{}
	AppID    string
	BrokerID string
	Event    string
	Time     int64
}
