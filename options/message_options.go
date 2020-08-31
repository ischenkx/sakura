package options

type MessageSend struct {
	Users []string
	Clients []string
	Channels []string
	Data []byte
	ToBeStored bool
	EventOptions
}