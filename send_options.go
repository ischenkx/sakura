package notify

type SendOptions struct {
	Users []string
	Clients []string
	Channels []string
	Message
}
