package notify

type MessageSendOptions struct {
	Users []string
	Clients []string
	Channels []string
	Data []byte
	ToBeStored bool
	EventOptions
}