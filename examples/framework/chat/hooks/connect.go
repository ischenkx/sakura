package hooks

import (
	"fmt"
	"github.com/ischenkx/notify"
)

// notify:hook name=connect priority = 123
func onConnect(app *notify.App, options notify.ConnectOptions, client notify.Client) {
	fmt.Println("connected")
}
