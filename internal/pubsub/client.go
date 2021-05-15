package pubsub

type Client interface {
	ID() string
	User() string
	Valid() bool
}
