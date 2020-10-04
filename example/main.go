package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/api/dns"
	"github.com/RomanIschenko/notify/auth/jwt"
	redibroker "github.com/RomanIschenko/notify/brokers/redis"
	"github.com/RomanIschenko/notify/events"
	nsj "github.com/RomanIschenko/notify/transports/sockjs"
	"github.com/RomanIschenko/pubsub"
	"github.com/RomanIschenko/pubsub/publication"
	"github.com/go-redis/redis/v8"
	"github.com/igm/sockjs-go/sockjs"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

func main() {

	PORT, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println(err)
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:               "redis-18197.c233.eu-west-1-1.ec2.cloud.redislabs.com:18197",
		Password:           "aByovWxGGvKQFRQG2DJfB7q0UzBMaeEf",
	})

	broker := redibroker.New(redisClient)

	go broker.Run(context.Background(), 16)

	server := nsj.NewServer("/pubsub", sockjs.DefaultOptions)
	http.HandleFunc("/pubsub/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		fmt.Println(r.URL.String())
		server.ServeHTTP(w, r)
	})

	app := notify.New(notify.Config{
		ID:               "app",
		Broker:           broker,
		PubSub:           pubsub.Config{},
		Server:           server,
		ServerGoroutines: 12,
		Auth:             jwt.New("thisIsTheJwtSecretPassword"),
		DataHandler:      func(app *notify.App, data notify.IncomingData) error {
								mes, err := ioutil.ReadAll(data.Reader)
								if err != nil {
									return err
								}
								app.Publish(pubsub.PublishOptions{
									Topics:  []string{"chat"},
									Payload: publication.New(mes),
								})
								return nil
							},
	})

	closer := app.Events().Handle(func(e events.Event) {
		if e.Type == notify.ConnectEvent {
			if client, ok := e.Data.(*pubsub.Client); ok {
				app.Subscribe(pubsub.SubscribeOptions{
					Topics:  []string{"chat"},
					Clients: []string{string(client.ID())},
				})
			}
		}
	})
	defer closer.Close()

	ip := fmt.Sprintf("localhost:%d", PORT)

	dns.HandleDNSPingEvent(broker, func() string {
		return ip
	})

	ctx, cancel := context.WithCancel(dns.WithIP(context.Background(), ip))

	wg := &sync.WaitGroup{}

	go app.Start(ctx, wg)

	app.Events().Handle(func(e events.Event) {
		fmt.Println("new event:", e.Type)
	})

	go http.ListenAndServe(ip, nil)

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt)

	select {
	case <-ctx.Done():
	case <-signals:
		cancel()
		wg.Wait()
	}
}
