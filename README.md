STILL IN DEVELOPMENT!!!

Simple echo server

```go
package main

import (
	"context"
	"github.com/ischenkx/swirl"
	authmock "github.com/ischenkx/swirl/pkg/auth/mock"
	"github.com/ischenkx/swirl/pkg/default/batchproto"
	evcodec "github.com/ischenkx/swirl/pkg/default/event_codec"
	"github.com/ischenkx/swirl/pkg/transports/websockets"
	"log"
	"net/http"
)

func main() {
	app := swirl.New(swirl.Config{
		EventsCodec:      evcodec.JSON{},
		ProtocolProvider: batchproto.NewProvider(1024),
		Auth:             authmock.New(),
		Adapter:          nil,
	})

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