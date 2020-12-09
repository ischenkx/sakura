package pubsub

import (
	"github.com/RomanIschenko/notify/pubsub/transport"
	"time"
)

type DisconnectOptions struct {
	Clients []string
	Users	[]string
	All bool
	Time int64
}

func (t *DisconnectOptions) validate() {
	if t.Time <= 0 {
		t.Time = time.Now().UnixNano()
	}
}

type PublishOptions struct {
	Topics []string
	Clients []string
	Users []string
	Payload []byte
	Time int64
	MetaInfo interface{}
}

func (t *PublishOptions) validate() {
	if t.Time <= 0 {
		t.Time = time.Now().UnixNano()
	}
}

type SubscribeOptions struct {
	Topics []string
	Clients []string
	Users []string
	Time int64
	MetaInfo interface{}
}

func (t *SubscribeOptions) validate() {
	if t.Time <= 0 {
		t.Time = time.Now().UnixNano()
	}
}

type UnsubscribeOptions struct {
	Topics []string
	Clients []string
	Users []string
	All bool
	Time int64
	MetaInfo interface{}
}

func (t *UnsubscribeOptions) validate() {
	if t.Time <= 0 {
		t.Time = time.Now().UnixNano()
	}
}

type ConnectOptions struct {
	Transport transport.Transport
	ID    	  string
	Time      int64
}
