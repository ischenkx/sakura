package startup

import (
	"context"
	"fmt"
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/pkg/transports/websockets"
	"net/http"
)

// notify:start
func Start(app *notify.App) {
	app.Start(context.Background())
	ws := websockets.NewServer(app.Server(), websockets.Config{})

	fmt.Println("Started!!!")

	http.ListenAndServe(":8080", ws)
}
