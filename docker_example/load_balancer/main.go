package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify/cluster/balancer"
	"github.com/RomanIschenko/notify/cluster/balancer/policies"
	redibroker "github.com/RomanIschenko/notify/cluster/broker/redis"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "host.docker.internal:6379",
	})

	adapter := redibroker.New(redisClient)

	b := balancer.New(balancer.Config{
		Policy:          &policies.RoundRobin{},
	})

	b.Start(context.Background(), adapter.Broker(ctx))

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		if addr, err := b.GetAddr(); err == nil {
			w.Write([]byte(addr))
		} else {
			w.WriteHeader(500)
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", 8888), nil)
}
