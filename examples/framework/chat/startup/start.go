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
	ws := websockets.NewServer(app.Server(), websockets.Config{})

	fmt.Println("Started!!!")

	http.ListenAndServe(":8080", ws)
}
