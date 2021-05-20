package simplehandler

import (
	"github.com/ischenkx/notify/framework/info"
	"log"
	"net/http"

	"github.com/ischenkx/notify/pkg/transports/websockets"
)

// notify:start
func Start(info info.Info) {
	app := info.App()
	app.Start(info.Context())

	log.Println("Listening websockets on port 8080...")

	ws := websockets.NewServer(info.App().Server(), websockets.Config{})

	if err := http.ListenAndServe(":8080", ws); err != nil {
		log.Println("Failed to serve:", err)
	}
}
