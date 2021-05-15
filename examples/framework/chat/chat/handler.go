package chat

import (
	"fmt"
	"github.com/RomanIschenko/notify"
)

// notify:handler
type Handler struct {
	app *notify.App
}

// notify:on event=message
func (h *Handler) HandleMessage() {
	fmt.Println("hello from sub-handler")
}
