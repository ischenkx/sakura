package main

import (
	"context"
	"encoding/json"
	"fmt"
	authmock "github.com/RomanIschenko/notify/auth/mock"
	"github.com/RomanIschenko/notify/cluster/api"
	redibroker "github.com/RomanIschenko/notify/cluster/broker/redis"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "host.docker.internal:6379",
	})

	adapter := redibroker.New(redisClient)

	if err != nil {
		panic(err)
		return
	}

	notifyApi := api.New(api.Config{
		AppID:  "app",
		Auth:   authmock.New(),
		Broker: adapter.Broker(context.Background()),
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var opts pubsub.PublishOptions

		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			logrus.Error(err)
			return
		}
		r.Body.Close()
		if err = notifyApi.Publish(opts); err != nil {
			logrus.Error(err)
		}
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		var opts pubsub.SubscribeOptions

		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			logrus.Error(err)
			return
		}
		r.Body.Close()
		if err = notifyApi.Subscribe(opts); err != nil {
			logrus.Error(err)
		}
	})

	http.HandleFunc("/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		var opts pubsub.UnsubscribeOptions

		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			logrus.Error(err)
			return
		}
		r.Body.Close()
		if err = notifyApi.Unsubscribe(opts); err != nil {
			logrus.Error(err)
		}
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Error(err)
			return
		}
		r.Body.Close()
		id, err := notifyApi.Authorize(string(body))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write([]byte(id))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}