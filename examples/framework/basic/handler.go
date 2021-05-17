package basic

import (
	"fmt"
	"github.com/ischenkx/notify"
	data2 "github.com/ischenkx/notify/examples/framework/basic/services/data"
	other2 "github.com/ischenkx/notify/examples/framework/basic/services/other"
)

// notify:handler
type Handler struct {
	// notify:inject
	app *notify.App

	// notify:inject
	s1 *data2.Service

	// notify:inject
	s2 *other2.Service
}

// notify:on event=message
func (h *Handler) HandleMessage() {
	fmt.Println("handling a message :", h.s1.OtherService)
}
