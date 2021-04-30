package pubsub

type Client interface {
	ID() string
	User() string
	Data() LocalStorage
	UserData() LocalStorage
}
