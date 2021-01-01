package main

import (
	"context"
	"github.com/RomanIschenko/notify"
	authmock "github.com/RomanIschenko/notify/auth/mock"
	"github.com/RomanIschenko/notify/pubsub"
	sjs "github.com/RomanIschenko/notify/transports/sockjs"
	events "github.com/RomanIschenko/notify/util/events"
	"github.com/igm/sockjs-go/v3/sockjs"
	"net/http"
	"time"
)

type message struct {
	Data string `json:"data"`
	Topics []string `json:"topics"`
}

func main() {
	sockjsOpts := sockjs.DefaultOptions
	sockjsOpts.CheckOrigin = func(request *http.Request) bool {
		return true
	}
	server := sjs.NewServer("/pubsub", sockjsOpts)

	app := notify.New(notify.Config{
		ID:            "xmas-mafia",
		PubsubConfig:  pubsub.Config{
			Workers:          6,
			ClientTTL:        time.Minute,
			ClientBufferSize: -1,
		},
		CleanInterval: time.Minute,
		Server:        notify.ServerConfig{
			Instance: server,
			Workers:  6,
		},
		Auth:          authmock.New(),
	})

	emitter := events.NewEmitter(context.Background(), app, events.JSONCodec{})

	emitter.On("subscribe", func(app *notify.App, emitter *events.Emitter, client pubsub.Client, topics []string) {
		app.Subscribe(pubsub.SubscribeOptions{
			Topics:   topics,
			Clients:  notify.IDs{client.ID()},
		})
	})

	emitter.On("unsubscribe", func(app *notify.App, emitter *events.Emitter, client pubsub.Client, topics []string) {
		app.Unsubscribe(pubsub.UnsubscribeOptions{
			Topics:   topics,
			Clients:  notify.IDs{client.ID()},
		})
	})

	emitter.On("publish", func(emitter *events.Emitter, client pubsub.Client, m message) {
		emitter.Emit(events.Event{
			Name: "message",
			Topics: m.Topics,
			Data: m.Data,
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

	http.ListenAndServe("localhost:4545", nil)
}
