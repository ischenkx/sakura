package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify/api"
	"github.com/RomanIschenko/notify/api/dns"
	"github.com/RomanIschenko/notify/auth/jwt"
	redibroker "github.com/RomanIschenko/notify/brokers/redis"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

func main()  {

	redisClient := redis.NewClient(&redis.Options{
		Addr:               "redis-18197.c233.eu-west-1-1.ec2.cloud.redislabs.com:18197",
		Password:           "aByovWxGGvKQFRQG2DJfB7q0UzBMaeEf",
	})

	broker := redibroker.New(redisClient)

	go broker.Run(context.Background(), 16)

	s := api.NewServer(api.Config{
		DNS: dns.Config{
			Broker:          broker,
			PingInterval:    time.Minute * 2,
			PongsBufferSize: 100,
		},
		Auth:   jwt.New("thisIsTheJwtSecretPassword"),
	})
	go s.Start(context.Background())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		fmt.Println(r.URL.String())
		s.ServeHTTP(w, r)
	})
	http.ListenAndServe(":8888", nil)
}
