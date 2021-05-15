package pubsub

import (
	message2 "github.com/RomanIschenko/notify/internal/pubsub/message"
	"io"
	"time"
)

type TimeStamp int64

type PublishOptions struct {
	Clients   []string
	Users     []string
	Topics    []string
	Data      []message2.Message
	TimeStamp int64
	Meta      interface{}
}

type SubscribeClientOptions struct {
	ID string
	Topic string
	TimeStamp int64
	Meta      interface{}
}

type SubscribeUserOptions struct {
	ID string
	Topic string
	TimeStamp int64
	Meta      interface{}
}

type UnsubscribeClientOptions struct {
	ID string
	Topic string
	TimeStamp  int64
	All bool
	Meta       interface{}
}

type UnsubscribeUserOptions struct {
	ID string
	Topic string
	TimeStamp  int64
	All bool
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

func (opts *SubscribeClientOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *UnsubscribeClientOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *SubscribeUserOptions) Validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *UnsubscribeUserOptions) Validate() {
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
