# Notify
Tiny Pub/Sub library. It is supposed to delivery notifications
to your clients via WebSockets/UDP/TCP/(any other transport)

# Example
Let's create a simple chat server
```go
// server.go
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

func main() {

	// websockets server based on gobwas/ws
	server := websockets.NewServer(ws.DefaultUpgrader, ws.DefaultHTTPUpgrader)

	// here we create a simple application
	// to send notifications (chat messages)
	app := notify.New(notify.Config{
		ID:           "chat",
		PubSubConfig: pubsub.Config{
			ClientConfig:  pubsub.ClientConfig{
				// Client's time to live
				// after the specified duration client will become invalid
				TTL:              time.Minute * 2,
			},
			CleanInterval: time.Minute,
		},
		ServerConfig: notify.ServerConfig{
			// our websockets server
			Server:          server,
			// Server interface provides three channels to read from:
			// Inactive - dropped connections
			// Incoming - incoming data from clients
			// Accept - incoming connections to be registered in our app
			// All these channels are being readen in different 
			// goroutines and Workers is a parameter that specifies 
			// the amount of reading goroutines for each channel
			Workers: 6,
			// DataHandler is a function that handles incoming data
			// func(*notify.App, notify.IncomingData)
			DataHandler:     DataHandler,
		},
		// Auth can be used for authentication of incoming connections
		Auth:         authmock.New(),
	})

	// we want to subscribe client to the "chat" topic
	// and send a "welcome" message to that topic
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
		// cors settings
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		server.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Println(err)
	}
}

func DataHandler(app *notify.App, data notify.IncomingData) {
	app.Publish(pubsub.PublishOptions{
		Topics:   []string{"chat"},
		Payload:  data.Payload,
	})
}
```

# Roadmap
- [x] Pub/Sub
- [x] WebSocket transport
- [ ] Long-Polling transport
- [ ] Clear resending mechanics
- [ ] Delivery guarantees (database?)

# Cluster Roadmap
- [x] Basic cluster
- [ ] Proper implementation of dns load-balancer
- [ ] Proper implementation of api service
- [ ] Broker implementation for RabbitMQ and Kafka
