package hooks

import (
	"fmt"
	"github.com/ischenkx/notify"
)

// notify:hook name="reconnect"
func ReconnectHook(app *notify.App, opts notify.ConnectOptions, c notify.Client) {
	fmt.Println("reconnected:", c.ID())
}
