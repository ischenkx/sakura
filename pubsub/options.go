package pubsub

import (
	"github.com/RomanIschenko/notify/pubsub/message"
	"io"
	"time"
)

type TimeStamp int64

type PublishOptions struct {
	Clients   []string
	Users     []string
	Topics    []string
	Data      []message.Message
	TimeStamp int64
	Meta      interface{}
}

type SubscribeOptions struct {
	Clients   []string
	Users     []string
	Topics    []string
	TimeStamp int64
	Meta      interface{}
}

type UnsubscribeOptions struct {
	Clients      []string
	Users        []string
	All          bool
	AllFromTopic bool
	AllClients bool
	Topics     []string
	TimeStamp  int64
	Meta       interface{}
}

type DisconnectOptions struct {
	Clients   []string
	Users     []string
	All       bool
	TimeStamp int64
	Meta      interface{}
}

type ConnectOptions struct {
	ClientID  string
	UserID    string
	Writer    io.WriteCloser
	TimeStamp int64
	Meta      interface{}
}

func (opts *PublishOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *SubscribeOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *UnsubscribeOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *DisconnectOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *ConnectOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}
