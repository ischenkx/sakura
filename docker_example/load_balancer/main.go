package main

import (
	"context"
	"fmt"
	redibroker "github.com/RomanIschenko/notify/brokers/redis"
	dnslb "github.com/RomanIschenko/notify/load_balancer/dns"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	fmt.Println("dns load balancer started")
	logrus.SetLevel(logrus.TraceLevel)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "host.docker.internal:6379",
	})

	broker := redibroker.New(redibroker.Config{
		Client: redisClient,
	})

	broker.Start(context.Background())

	server := dnslb.NewServer(dnslb.Config{
		Broker:          broker,
		PingInterval: time.Second * 10,
	})

	lis, err := net.Listen("tcp", ":8888")

	if err != nil {
		panic(err)
	}

	server.Start(context.Background(), lis)
}
