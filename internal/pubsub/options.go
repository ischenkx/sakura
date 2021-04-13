package pubsub

import (
	"io"
	"time"
)

type TimeStamp int64

type PublishOptions struct {
	Clients   []string
	Users     []string
	Topics    []string
	Data      []byte
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
	// NOT IMPLEMENTED!
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

func (opts *PublishOptions) validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *SubscribeOptions) validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *UnsubscribeOptions) validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *DisconnectOptions) validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}

func (opts *ConnectOptions) validate() {
	if opts.TimeStamp <= 0 {
		opts.TimeStamp = time.Now().UnixNano()
	}
}
