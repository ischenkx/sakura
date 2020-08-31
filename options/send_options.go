package options

import "github.com/RomanIschenko/notify"

type Send struct {
	Users 	   []string
	Clients    []string
	Channels   []string
	ToBeStored bool
	notify.Message
	EventOptions
}
