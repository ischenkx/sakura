package pubsub

import "io"

type SubscribeOptions struct {
	Topics []string
	Clients, Users []string
	MetaInfo interface{}
	Seq int64
}

type UnsubscribeOptions struct {
	Topics []string
	Clients, Users []string
	MetaInfo interface{}
	Seq int64
}

type PublishOptions struct {
	Topics []string
	Clients, Users []string
	Message []byte
	NoBuffering bool
	MetaInfo interface{}
	Seq int64
}

type DisconnectOptions struct {
	Clients, Users []string
	All bool
	MetaInfo interface{}
	Seq int64
}

type InactivateOptions struct {
	Time int64
	ClientID string
}

type ConnectOptions struct {
	ClientID, UserID string
	Transport io.WriteCloser
	MetaInfo interface{}
	Seq int64
}