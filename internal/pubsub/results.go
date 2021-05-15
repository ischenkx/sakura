package pubsub

type ConnectResult struct {
	ChangeLog   *ChangeLog
	Error       error
	Client      Client
	Reconnected bool
}