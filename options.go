package swirl

type SubscribeOptions struct {
	Topics    []string
	TimeStamp int64
	MetaInfo  interface{}
}

type UnsubscribeOptions struct {
	Topics    []string
	All       bool
	TimeStamp int64
	MetaInfo  interface{}
}

type EventOptions struct {
	Name      string
	Args      []interface{}
	TimeStamp int64
	MetaInfo  interface{}
}

type DisconnectOptions struct {
	TimeStamp int64
	MetaInfo  interface{}
}
