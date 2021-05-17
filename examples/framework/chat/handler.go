package chat

import (
	"github.com/ischenkx/notify"
)

// notify:handler
type Chat struct {
	// notify:inject
	app *notify.App

	// notify:inject
	service *Service
}

// notify:on event=message
func (c *Chat) HandleMessage(client notify.Client) {
	client.Emit("message", "hello from message handler:" + c.service.Name)
}


// notify:handler prefix="handler0"
type Handler0 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler0) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 0 handler")
}



// notify:handler prefix="handler1"
type Handler1 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler1) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 1 handler")
}



// notify:handler prefix="handler2"
type Handler2 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler2) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 2 handler")
}



// notify:handler prefix="handler3"
type Handler3 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler3) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 3 handler")
}

// notify:handler prefix="handler4"
type Handler4 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler4) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 4 handler")
}



// notify:handler prefix="handler5"
type Handler5 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler5) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 5 handler")
}

// notify:handler prefix="handler6"
type Handler6 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler6) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 6 handler")
}



// notify:handler prefix="handler7"
type Handler7 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler7) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 7 handler")
}



// notify:handler prefix="handler8"
type Handler8 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler8) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 8 handler")
}



// notify:handler prefix="handler9"
type Handler9 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler9) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 9 handler")
}



// notify:handler prefix="handler10"
type Handler10 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler10) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 10 handler")
}



// notify:handler prefix="handler11"
type Handler11 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler11) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 11 handler")
}



// notify:handler prefix="handler12"
type Handler12 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler12) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 12 handler")
}



// notify:handler prefix="handler13"
type Handler13 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler13) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 13 handler")
}



// notify:handler prefix="handler14"
type Handler14 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler14) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 14 handler")
}



// notify:handler prefix="handler15"
type Handler15 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler15) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 15 handler")
}



// notify:handler prefix="handler16"
type Handler16 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler16) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 16 handler")
}



// notify:handler prefix="handler17"
type Handler17 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler17) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 17 handler")
}



// notify:handler prefix="handler18"
type Handler18 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler18) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 18 handler")
}



// notify:handler prefix="handler19"
type Handler19 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler19) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 19 handler")
}



// notify:handler prefix="handler20"
type Handler20 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler20) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 20 handler")
}



// notify:handler prefix="handler21"
type Handler21 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler21) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 21 handler")
}



// notify:handler prefix="handler22"
type Handler22 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler22) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 22 handler")
}



// notify:handler prefix="handler23"
type Handler23 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler23) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 23 handler")
}



// notify:handler prefix="handler24"
type Handler24 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler24) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 24 handler")
}



// notify:handler prefix="handler25"
type Handler25 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler25) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 25 handler")
}



// notify:handler prefix="handler26"
type Handler26 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler26) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 26 handler")
}



// notify:handler prefix="handler27"
type Handler27 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler27) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 27 handler")
}



// notify:handler prefix="handler28"
type Handler28 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler28) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 28 handler")
}



// notify:handler prefix="handler29"
type Handler29 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler29) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 29 handler")
}



// notify:handler prefix="handler30"
type Handler30 struct {
	// notify:inject
	app *notify.App
}

// notify:on event=message
func (h *Handler30) HandleMessage(c notify.Client) {
	c.Emit("message", "response from 30 handler")
}
