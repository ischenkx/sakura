package sockjs

import (
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var logger = logrus.WithField("source", "sockjs_server")

type Server struct {
	mu      sync.RWMutex
	app     notify.Servable
	handler http.Handler
}

func (s *Server) serveSockJS(session sockjs.Session) {
	t := &Transport{
		session: session,
	}
	auth, err := t.session.Recv()
	if err != nil {
		logger.Errorf("failed to recv data, closing:", err)
		t.Close()
		return
	}

	opts := pubsub.ConnectOptions{
		ClientID:  "",
		UserID:    "",
		Writer:    nil,
		TimeStamp: 0,
		Meta:      nil,
	}
	s.mu.RLock()
	app := s.app
	s.mu.RUnlock()
	client, err := app.Connect(auth, opts)
	if err != nil {
		fmt.Println("connection failure:", err)
		t.Close()
		return
	}
	for {
		data, err := t.session.Recv()
		if err != nil {
			app.Inactivate(client.ID())
			return
		}
		app.HandleMessage(client, []byte(data))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func NewServer(app notify.Servable, prefix string, opts sockjs.Options) *Server {
	s := &Server{}
	s.app = app
	s.handler = sockjs.NewHandler(prefix, opts, s.serveSockJS)
	return s
}
