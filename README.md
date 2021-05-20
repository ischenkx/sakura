# Notify

A tiny library for building real-time applications

# General information

The main object of your program is application. It contains
all information about clients, users and topics.

Client is a wrapper around some Transport (it can be
a websocket or a tcp connection).

User is simply a group of clients.

Both of structures described above can be subscribed to topics - abstract group of clients (for example, a chat room)

# Examples

The following examples uses code generation based on high
level meta-information (via comments starting with notify:)

```go
package app

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/framework/builder"
	"github.com/ischenkx/notify/framework/info"
	"github.com/ischenkx/notify/framework/ioc"
	authmock "github.com/ischenkx/notify/pkg/auth/mock"
	"github.com/ischenkx/notify/pkg/transports/websockets"
	"net/http"
)

type UniqueKeyProvider struct{}

func (c *UniqueKeyProvider) GenerateKey() string {
	return uuid.New().String()
}


// notify:configure
func Configure(b *builder.Builder) {
	fmt.Println("Configuring app...")
	b.AppConfig().ID = "my-first-app"
	b.AppConfig().Auth = authmock.New()
	b.IocContainer().AddEntry(ioc.NewEntry(&UniqueKeyProvider{}, ""))
}

// notify:initialize
func Initialize(info info.Info) {
	fmt.Println("Initializing app...")
}

// notify:start
func Start(info info.Info) {
	info.App().Start(info.Context())
	wsServer := websockets.NewServer(info.App().Server(), websockets.Config{})
	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", wsServer); err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}

// notify:hook name=connect
func ConnectHook(app *notify.App, opts notify.ConnectOptions, c notify.Client) {
	var userID = "<no user>"
	if u, err := c.User(); err == nil {
		userID = u.ID()
	}
	fmt.Printf("CONNECT_HOOK [id=%s user=%s]\n", c.ID(), userID)
}

// notify:handler
type Handler struct {
	// notify:inject
	app               *notify.App

	// notify:inject
	uniqueKeyProvider *UniqueKeyProvider
}

type Response struct {
	Message       string `json:"message,omitempty"`
	TransactionID string `json:"tid,omitempty"`
}

// notify:on event=message
func (h *Handler) OnMessage(c notify.Client, data string) {
	txID := h.uniqueKeyProvider.GenerateKey()
	c.Emit("message", Response{
		Message: data,
		TransactionID: txID,
	})
}
```
