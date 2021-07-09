package swirl

type Metrics interface {
	Clients() IDList
	Topics() IDList
	Users() IDList
}
