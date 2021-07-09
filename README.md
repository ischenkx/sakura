### About

A library for real-time communication in modern web applications

General concepts:
- Client - a single connection (for example a websocket) that's used for transferring events
- User - an abstract entity that is mainly used for storing persistent subscriptions
- Topic - an abstract group that users and clients can join or leave (for example some chat room)

#### Install
```
go get github.com/ischenkx/swirl
```

#### Example
Simple echo server
```go
package main

import (
	"context"
	"github.com/ischenkx/swirl"
	"github.com/ischenkx/swirl/pkg/transports/websockets"
	"log"
	"net/http"
)

func main() {
	app := swirl.New(swirl.Config{})

	app.Start(context.Background())

	app.On("echo", func(client swirl.Client, data string) {
		client.Emit("echo", swirl.Args{data})
	})

	server := websockets.NewServer(app.Server(), websockets.Config{})

	if err := http.ListenAndServe(":6060", server); err != nil {
		log.Fatalln(err)
	}
}

```