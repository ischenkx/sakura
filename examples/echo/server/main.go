package main

import (
	"context"
	"github.com/ischenkx/swirl"
	"github.com/ischenkx/swirl/pkg/transports/websockets"
	"log"
	"net/http"
)

func main() {

	// use default config
	app := swirl.New(swirl.Config{})

	app.On("echo", func(c swirl.Client, data interface{}) {
		c.Emit("echo", swirl.Args{data})
	})

	app.Events(swirl.UserPriority).OnError(func(app *swirl.App, err error) {
		log.Println("error:", err)
	})

	app.Start(context.Background())

	http.Handle("/swirl", websockets.NewServer(app.Server(), websockets.Config{}))
	log.Println("starting an http server...")
	if err := http.ListenAndServe(":6060", nil); err != nil {
		log.Fatal("failed to start an http server:", err)
	}
}
