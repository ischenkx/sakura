package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify/api"
	"github.com/RomanIschenko/notify/auth/jwt"
	redibroker "github.com/RomanIschenko/notify/brokers/redis"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net"
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

	broker := redibroker.New(redibroker.Config{
		Client:                 redisClient,
	})

	broker.Start(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		panic(err)
		return
	}

	server := api.NewServer(api.Config{
		APIKey:           "api_key",
		Auth:             jwt.New("secret_key"),
		Broker:           broker,
	})

	server.Start(context.Background(), lis)
}