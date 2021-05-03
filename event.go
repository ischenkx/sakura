package notify

type Event struct {
	Name      string
	Data      interface{}
	Clients   []string
	Users     []string
	Topics    []string
	Timestamp int64
	Meta      interface{}
}
