package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RomanIschenko/notify/api"
	dnslb "github.com/RomanIschenko/notify/load_balancer/dns"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"net/http"
)

type Output struct {
	IP, Token string
}

func main() {
	psConn, err := grpc.Dial("localhost:7878", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	lbConn, err := grpc.Dial("localhost:6868", grpc.WithInsecure())

	pubsubClient := api.NewClient("app", psConn)
	lbClient := dnslb.NewClient(lbConn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		id := uuid.New().String()
		authResult, err := pubsubClient.Authorize(context.Background(), pubsub.ClientID(id))

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		lbResult, err := lbClient.GetAddress(context.Background())

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		json.NewEncoder(w).Encode(Output{lbResult, authResult.Token})
	})

	fmt.Println(http.ListenAndServe("localhost:8888", nil))
}