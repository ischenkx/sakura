package sockjs

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

var logger = logrus.WithField("source", "sockjs_server")

type Server struct {
	mu sync.RWMutex
	app			 notify.ServableApp
	handler		 http.Handler
}

func (s *Server) serveSockJS(session sockjs.Session) {
	t := &Transport{
		session:   session,
	}
	auth, err := t.session.Recv()
	if err != nil {
		logger.Errorf("failed to recv data, closing:", err)
		t.Close()
		return
	}

	opts := pubsub.ConnectOptions{
			Transport: t,
			Seq:      time.Now().UnixNano(),
			MetaInfo:  session.Request(),
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
		app.HandleIncomingData(notify.IncomingData{
			Client:  client,
			Payload: []byte(data),
		})
	}
}


func (s *Server) Start(ctx context.Context, app notify.ServableApp) {
	s.mu.Lock()
	s.app = app
	s.mu.Unlock()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func NewServer(prefix string, opts sockjs.Options) *Server {
	s := &Server{}
	s.handler = sockjs.NewHandler(prefix, opts, s.serveSockJS)
	return s
}