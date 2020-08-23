package notify

type Message struct {
	Data []byte
	ID string
}

type MessageOptions struct {
	Users []string
	Clients []string
	Channels []string
	Data []byte
}
