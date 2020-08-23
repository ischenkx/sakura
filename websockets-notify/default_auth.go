package websockets_notify

import (
	"encoding/json"
	"net/http"
)

type DefaultAuthenticator struct {}

func (auth DefaultAuthenticator) Authenticate(r *http.Request) (ClientData, error) {
	creds := r.URL.Query().Get("auth")
	data :=  &ClientData{}
	err := json.Unmarshal([]byte(creds), data)
	return *data, err
}
