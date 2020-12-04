package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	authmock "github.com/RomanIschenko/notify/auth/mock"
	"github.com/RomanIschenko/notify/cluster"
	"github.com/RomanIschenko/notify/cluster/balancer"
	redibroker "github.com/RomanIschenko/notify/cluster/broker/redis"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	taskctx "github.com/RomanIschenko/notify/task_context"
	"github.com/RomanIschenko/notify/transports/websockets"
	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func handleData(app *notify.App, data notify.IncomingData) {
	app.Publish(pubsub.PublishOptions{
		Topics:  []string{"chat"},
		Clients: nil,
		Users:   nil,
		Payload: data.Payload,
	})
}

func main() {

	ctx, shutdown := context.WithCancel(context.Background())
	taskContext := taskctx.New(ctx)

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

	adapter := redibroker.New(redisClient)

	adapter.Broker(context.Background())

	server := websockets.NewServer(ws.DefaultUpgrader, ws.DefaultHTTPUpgrader)

	app := notify.New(notify.Config{
		ID:           "app",
		PubSubConfig: pubsub.Config{
			Shards:         16,
			ClientConfig: pubsub.ClientConfig{
				TTL:              time.Second * 10,
				InvalidationTime: time.Second * 10,
				BufferSize: 		1024,
			},
			CleanInterval:  time.Second*10,
		},
		ServerConfig: notify.ServerConfig{
			Server: server,
			DataHandler: handleData,
		},
		Auth:  authmock.New(),
	})
	app.Events(taskContext).OnConnect(func(opts pubsub.ConnectOptions, client *pubsub.Client, log changelog.Log) {
		app.Subscribe(pubsub.SubscribeOptions{
			Topics:  []string{"chat"},
			Clients: []string{client.ID().String()},
		})
		app.Publish(pubsub.PublishOptions{
			Topics:  []string{"chat"},
			Payload: []byte(fmt.Sprintf("client %s joined the chat %s", client.ID(), "chat")),
		})
	})

	err = cluster.RegisterApp(ctx, adapter.Broker(ctx), app)
	if err != nil {
		panic(err)
	}
	err = balancer.Register(taskContext, adapter.Broker(ctx), app, func() string {
		return fmt.Sprintf("localhost:%d", outterPort)
	})

	if err != nil {
		panic(err)
	}

	app.Start(ctx)

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
	signal.Notify(intChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(time.Second * 15)
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
	}
	shutdown()
	logrus.WithField("source", "main").Debug("EXIT")
}