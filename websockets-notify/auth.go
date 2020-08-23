package websockets_notify

import "net/http"

type ClientData struct {
	ID 	   string
	UserID string
	App	   string
}

type Authenticator interface {
	Authenticate(r *http.Request) (ClientData, error)
}
