package notify

type BrokerEvent struct {
	Data     interface{}
	AppID    string
	BrokerID string
	Event    string
	Time     int64
}
