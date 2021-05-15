package hooks

import (
	"fmt"
	"github.com/RomanIschenko/notify"
)

// notify:hook name="beforeConnect"
func OnBeforeConnect(app *notify.App, opts *notify.ConnectOptions) {
	fmt.Println("before connect:", opts.ClientID)
}