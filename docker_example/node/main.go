package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/auth/jwt"
	redibroker "github.com/RomanIschenko/notify/brokers/redis"
	"github.com/RomanIschenko/notify/events"
	dnslb "github.com/RomanIschenko/notify/load_balancer/dns"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/publication"
	"github.com/RomanIschenko/notify/transports/websockets"
	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func handleData(app *notify.App, data notify.IncomingData) error {
	bts, err := ioutil.ReadAll(data.Reader)
	if err != nil {
		return err
	}
	
	app.Publish(pubsub.PublishOptions{
		Topics:  []string{"chat"},
		Clients: nil,
		Users:   nil,
		Payload: publication.New(bts),
	})

	return nil
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	innerPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

    res, err := http.Get("http://reg:9090/port?id="+hostname)
	if err != nil {
		panic(err)
	}

	portBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	outterPort, err := strconv.Atoi(string(portBytes))
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "host.docker.internal:6379",
	})

	broker := redibroker.New(redibroker.Config{
		Client: redisClient,
	})

	broker.Start(context.Background())

	server := websockets.NewServer(ws.DefaultUpgrader, ws.DefaultHTTPUpgrader)

	app := notify.New(notify.Config{
		ID:           "app",
		Broker:       broker,
		PubSubConfig: pubsub.Config{
			Shards:         16,
			ShardConfig: pubsub.ShardConfig{
				ClientTTL:              time.Second * 10,
				ClientInvalidationTime: time.Second * 10,
			},
			CleanInterval:  time.Second*10,
			TopicBuckets:   6,
		},
		ServerConfig: notify.ServerConfig{
			Server: server,
			Goroutines: 64,
			DataHandler: handleData,
		},
		Auth:         jwt.New("secret_key"),
	})
	
	appEventsHandle := app.Events().Handle(func(event events.Event) {
		if event.Type == notify.ConnectEvent {
			client := event.Data.(*pubsub.Client)

			app.Subscribe(pubsub.SubscribeOptions{
				Topics:  []string{"chat"},
				Clients: []string{client.ID().String()},
			})

			app.Publish(pubsub.PublishOptions{
				Topics:  []string{"chat"},
				Payload: publication.New(
					[]byte(fmt.Sprintf("client %s joined the chat", client.ID())),
				),
			})
		}
	})
	defer appEventsHandle.Close()

	ctx, cancel := context.WithCancel(context.Background())

	awaiter := app.Start(ctx)

	regHandle, err := dnslb.Register(app, broker, fmt.Sprintf("localhost:%d", outterPort))

	if err != nil {
		panic(err)
	}

	defer regHandle.Close()

	http.HandleFunc("/pubsub", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		server.ServeHTTP(w, r)
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", innerPort), nil)
	
	intChan := make(chan os.Signal)
	signal.Notify(intChan, os.Interrupt)

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				metrics := app.Metrics()

				logrus.WithField("source", "metrics").
					Debugf("Users: %d\nClients: %d\nTopics: %d", metrics.Users, metrics.Clients, metrics.Topics)
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case <- ctx.Done():
	case <- intChan:
		cancel()
	}

	awaiter.Wait()
	logrus.WithField("source", "main").Debug("finishing")
}