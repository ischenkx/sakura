package data

import (
	other2 "github.com/ischenkx/notify/examples/framework/basic/services/other"
)

type Service struct {
	Data string
	// notify:inject
	OtherService *other2.Service
}
