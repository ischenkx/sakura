package parser

import "github.com/ischenkx/notify"

type Hook struct {
	Path string
	Name string
	FuncName string
	Priority notify.Priority
}