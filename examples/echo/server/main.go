package main

import (
	"context"
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/pkg/auth/jwt"
	"github.com/ischenkx/notify/pkg/default/codec"
	"github.com/ischenkx/notify/pkg/transports/websockets"
	"github.com/gobwas/ws"
	"log"
	"net/http"
)

const JwtSecret = "javainuse-secret-key"

func main() {
	app := notify.New(notify.Config{
		ID:    "echo-server",
		Auth:  jwt.New(JwtSecret),
		Codec: codec.JSON{},
	})
	hooks := app.Hooks(notify.UserPriority)

	hooks.OnEvent(func(app *notify.App, client notify.Client, event notify.IncomingEvent) {
		//log.Printf("EVENT [name]=%v [data]=%v\n", event.Name, event.State)
	})

	hooks.OnFailedEvent(func(app *notify.App, name string, payload []byte) {
		log.Printf("FAILED EVENT name=[%s]\n payload=[%-20s]", name, payload)
	})

	hooks.OnConnect(func(app *notify.App, options notify.ConnectOptions, client notify.Client) {
		userID := "<empty>"
		if user, err := client.User(); err == nil {
			userID = user.ID()
		}
		log.Printf("CONNECT id=[%s] user=[%s]\n", client.ID(), userID)
	})

	app.On("message", func(c notify.Client, incoming string) {
		c.Emit("message", incoming)
	})

	app.Start(context.Background())

	wsServer := websockets.NewServer(app.Server(), websockets.Config{
		Upgrader:        ws.DefaultHTTPUpgrader,
	})
	http.ListenAndServe("localhost:4545", wsServer)
}
