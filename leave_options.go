package notify

type LeaveOptions struct {
	Users []string
	Clients []string
	Channels []string
	All bool
	Event string
}