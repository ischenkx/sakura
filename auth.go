package notify

type Auth interface {
	Verify(interface{}) (ClientInfo, bool)
}
