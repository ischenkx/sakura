package startup

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pkg/transports/websockets"
	"net/http"
)

// notify:starter
func Start(app *notify.App) {
	app.Start(context.Background())
	server := websockets.NewServer(app.Server(), websockets.Config{})
	fmt.Println("running websocket server...")
	http.ListenAndServe(":8080", server)
}
