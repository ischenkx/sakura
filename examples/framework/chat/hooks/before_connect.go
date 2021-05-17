package hooks

import (
	"fmt"
	"github.com/ischenkx/notify"
)

// notify:hook name=beforeConnect priority=plugin
func OnBeforeConnect(app *notify.App, opts *notify.ConnectOptions) {
	fmt.Println("before connect:", opts.ClientID)
}