
package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	authmock "github.com/RomanIschenko/notify/auth/mock"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/RomanIschenko/notify/transports/websockets"
	"github.com/gobwas/ws"
	"log"
	"net/http"
	"time"
)

func DataHandler(app *notify.App, data notify.IncomingData) {
	app.Publish(pubsub.PublishOptions{
		Topics:   []string{"chat"},
		Payload:  data.Payload,
	})
}

func main() {

	server := websockets.NewServer(ws.DefaultUpgrader, ws.DefaultHTTPUpgrader)

	app := notify.New(notify.Config{
		ID:           "app",
		PubSubConfig: pubsub.Config{
			Shards:        16,
			ClientConfig:  pubsub.ClientConfig{
				TTL:              time.Minute * 2,
				InvalidationTime: time.Minute * 1,
				BufferSize:       1024,
			},
			CleanInterval: time.Minute,
		},
		ServerConfig: notify.ServerConfig{
			Server:          server,
			Workers: 6,
			DataHandler:     DataHandler,
		},
		Auth:         authmock.New(),
	})

	app.Events(context.Background()).
		OnConnect(func(opts pubsub.ConnectOptions, client *pubsub.Client, log changelog.Log) {
			app.Subscribe(pubsub.SubscribeOptions{
				Topics:   []string{"chat"},
				Clients:  []string{client.ID()},
			})
			app.Publish(pubsub.PublishOptions{
				Topics:   []string{"chat"},
				Payload:  []byte(fmt.Sprintf("%s joined the chat!!!", client.ID())),
			})
		})

	app.Start(context.Background())

	http.HandleFunc("/pubsub", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		server.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Println(err)
	}
}

