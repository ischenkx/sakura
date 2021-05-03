package main

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	authmock "github.com/RomanIschenko/notify/pkg/auth/mock"
	"github.com/RomanIschenko/notify/pkg/transports/websockets"
	"github.com/gobwas/ws"
	"log"
	"net/http"
)

type ChatJoinRequest struct {
	Name string
}

type ChatError struct {
	Message string
}

type ChatMessage struct {
	// 100 - just a simple message
	//
	// 200 - a system message
	Code int
	From string
	Payload string
}

func main() {
	// As it is an example let's use a mock auth
	app := notify.New(notify.Config{
		ID:            "chat-app",
		Auth:          authmock.Auth{},
	})

	// ok, now let's handle some events
	hooks := app.Hooks(notify.UserPriority)
	defer hooks.Close()

	hooks.OnConnect(func(app *notify.App, opts notify.ConnectOptions, client notify.Client) {
		log.Println("connect")
	})

	hooks.OnDisconnect(func(app *notify.App, client notify.Client) {
		log.Println("disconnect")
	})

	hooks.OnInactivate(func(app *notify.App, client notify.Client) {
		log.Println("inactivate")
	})

	hooks.OnReconnect(func(app *notify.App, opts notify.ConnectOptions, client notify.Client) {
		log.Println("reconnect")
	})

	hooks.OnChange(func(app *notify.App, l notify.ChangeLog) {
		log.Println("change")
	})

	hooks.OnEmit(func(app *notify.App, event notify.Event) {
		log.Println("emit")
	})

	hooks.OnEvent(func(app *notify.App, client notify.Client, event notify.IncomingEvent) {
		log.Println("event")
	})

	hooks.OnBeforeConnect(func(app *notify.App, options *notify.ConnectOptions) {
		log.Println("before connect")
	})

	hooks.OnBeforeDisconnect(func(app *notify.App, options *notify.DisconnectOptions) {
		log.Println("before disconnect")
	})

	hooks.OnBeforeEmit(func(app *notify.App, event *notify.Event) {
		log.Println("before emit")
	})

	hooks.OnBeforeSubscribe(func(app *notify.App, options *notify.SubscribeOptions) {
		log.Println("before subscribe")
	})

	hooks.OnBeforeUnsubscribe(func(app *notify.App, options *notify.UnsubscribeOptions) {
		log.Println("before unsubscribe")
	})

	app.On("chat.join", func(c notify.Client, app *notify.App, request ChatJoinRequest) {
		if previousChat, ok := c.Data().LoadAndDelete("current_chat"); ok {
			app.Action().
				WithClients(c.ID()).
				WithTopics(previousChat.(string)).
				Unsubscribe()

			app.Emit(notify.Event{
				Name: "chat.message",
				Topics: notify.IDs{previousChat.(string)},
				Data: ChatMessage{
					Code: 200,
					Payload: fmt.Sprintf("%s left the chat", c.ID()),
				},
			})
		}

		if _, ok := c.Data().LoadOrStore("current_chat", request.Name); !ok {
			res := app.Action().
				WithClients(c.ID()).
				WithTopics(request.Name).
				Subscribe()

			name := c.ID()

			if n, ok := c.Data().Load("name"); ok {
				name = n.(string)
			}

			app.Emit(notify.Event{
				Name:      "chat.message",
				Data:      ChatMessage{
					Code:    200,
					From:    app.ID(),
					Payload: fmt.Sprintf("%s joined the chat", name),
				},
				Topics:    notify.IDs{request.Name},
			})

			log.Println("subscribed!!!", res.OK())
		}
	})

	app.On("chat.leave", func(c notify.Client, app *notify.App) {
		if chat, ok := c.Data().LoadAndDelete("current_chat"); ok {
			res := app.Action().
				WithClients(c.ID()).
				WithTopics(chat.(string)).
				Unsubscribe()

			if res.OK() {
				app.Emit(notify.Event{
					Name: "chat.message",
					Topics: notify.IDs{chat.(string)},
					Data: ChatMessage{
						Code: 200,
						Payload: fmt.Sprintf("%s left the chat", c.ID()),
					},
				})
			}
		}
	})

	app.On("chat.message", func(client notify.Client, app *notify.App, message ChatMessage) {
		if message.Code != 100 {
			app.Emit(notify.Event{
				Name:      "chat.error",
				Data:      ChatMessage{
					Code:    300,
					From: 	 app.ID(),
					Payload: fmt.Sprintf("invalid message code: %d", message.Code),
				},
				Clients:   notify.IDs{client.ID()},
			})
			return
		}

		if chat, ok := client.Data().Load("current_chat"); ok {
			if name, ok := client.Data().Load("name"); ok {
				message.From = name.(string)
			} else {
				message.From = client.ID()
			}
			app.Emit(notify.Event{
				Name:      "chat.message",
				Data:      message,
				Topics:    notify.IDs{chat.(string)},
			})
		}
	})

	app.On("user.set_name", func(c notify.Client, name string) {
		c.Data().Store("name", name)
	})

	app.Start(context.Background())

	server := websockets.NewServer(app.Servable(), ws.DefaultHTTPUpgrader)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})

	http.ListenAndServe("localhost:3434", nil)
}





