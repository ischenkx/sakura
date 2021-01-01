package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/auth/jwt"
	"github.com/RomanIschenko/notify/pubsub"
	sjs "github.com/RomanIschenko/notify/transports/sockjs"
	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	opts := sockjs.DefaultOptions
	opts.CheckOrigin = func(request *http.Request) bool {
		return true
	}
	opts.RawWebsocket = true
	server := sjs.NewServer("/pubsub", opts)

	app := notify.New(notify.Config{
		ID: "chat",
		PubsubConfig: pubsub.Config{
			ClientTTL: time.Second * 30,
			ClientBufferSize: -1,
		},
		Server: notify.ServerConfig{
			Instance:  server,
			Workers: 6,
		},
		CleanInterval: time.Second * 15,
		Auth: jwt.New("secret-key"),
	})

	app.Events(context.Background()).
		OnPublish(func(app *notify.App, options pubsub.PublishOptions) {
			//fmt.Println("publishing:", string(options.Message))
		})

	app.Events(context.Background()).
		OnConnect(func(app *notify.App, opts pubsub.ConnectOptions, client pubsub.Client) {
			if req, ok := opts.MetaInfo.(*http.Request); ok {
				room := req.URL.Query().Get("room")
				if room == "" {
					room = "chat"
				}
				client.Meta().Store("room", room)
				app.Subscribe(pubsub.SubscribeOptions{
					Topics:   []string{room},
					Clients:  []string{client.ID()},
				})
				app.Publish(pubsub.PublishOptions{
					Topics:  []string{room},
					Message: []byte(fmt.Sprintf("%s joined the room %s!!!", client.ID(), room)),
				})
			} else {
				fmt.Println("failed to get http request from options")
			}
		})

	app.Events(context.Background()).
		OnIncomingData(
			func(app *notify.App, data notify.IncomingData) {
				room, ok := data.Client.Meta().Load("room")
				if !ok {
					return
				}
				app.Publish(pubsub.PublishOptions{
					Topics:  []string{room.(string)},
					Message: data.Payload,
				})
			})

	app.Start(context.Background())

	http.HandleFunc("/pubsub/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		server.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Println(err)
	}
}