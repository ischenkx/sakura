package main

import (
	"context"
	"github.com/RomanIschenko/notify"
	authmock "github.com/RomanIschenko/notify/auth/mock"
	"github.com/RomanIschenko/notify/events"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/publication"
	nsj "github.com/RomanIschenko/notify/transports/sockjs"
	"github.com/igm/sockjs-go/sockjs"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	server := nsj.NewServer("/pubsub", sockjs.DefaultOptions)
	app := notify.New(notify.Config{
		ID:               "app",
		PubSubConfig:           pubsub.Config{
			Shards:         16,
			ShardConfig:    pubsub.ShardConfig{
				ClientTTL:              time.Minute*2,
				ClientInvalidationTime: time.Minute*2,
				ClientBufferSize:       500,
			},
			PubQueueConfig: pubsub.PubQueueConfig{
				BufferSize:       100,
				Writers:          12,
				ReadersPerWriter: 12,
			},
			CleanInterval:  time.Minute*4,
			TopicBuckets:   16,
		},
		ServerConfig: notify.ServerConfig{
			Server:      server,
			Goroutines:  10,
			DataHandler: func(app *notify.App, data notify.IncomingData) error {
				dataBytes, err := ioutil.ReadAll(data.Reader)
				if err != nil {
					return err
				}
				app.Publish(pubsub.PublishOptions{
					Topics:  []string{"chat"},
					Payload: publication.New(dataBytes),
				})
				return nil
			},
		},
		Auth: authmock.New(),
	})

	handle := app.Events().Handle(func(e events.Event) {
		if e.Type == notify.ConnectEvent {
			client := e.Data.(*pubsub.Client)
			app.Subscribe(pubsub.SubscribeOptions{
				Topics:  []string{"chat"},
				Clients: []string{client.ID().String()},
			})
		}
	})
	defer handle.Close()

	http.HandleFunc("/pubsub/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		server.ServeHTTP(w, r)
	})
	go http.ListenAndServe("localhost:6565", nil)
	app.Start(context.Background()).Wait()
}
