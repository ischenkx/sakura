package options

type Leave struct {
	Users []string
	Clients []string
	Channels []string
	All bool
	EventOptions
}