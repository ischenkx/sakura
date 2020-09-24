package sockjs

import (
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/pubsub"
	"github.com/igm/sockjs-go/sockjs"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	acceptChan chan notify.IncomingConnection
	incomingChan chan notify.IncomingData
	inactivateChan chan *pubsub.Client
	handler		 http.Handler
}

func (s *Server) serveSockJS(session sockjs.Session) {
	t := &Transport{
		session:   session,
		state:  int32(pubsub.OpenTransport),
	}

	auth, err := t.session.Recv()

	if err != nil {
		t.Close()
		return
	}

	conn := notify.IncomingConnection{
		Opts: pubsub.ConnectOptions{
			Transport: t,
			Time:      time.Now().UnixNano(),
		},
		AuthData: auth,
		Resolver: make(chan notify.ResolvedConnection),
	}

	s.acceptChan <- conn

	resolvedConn := <-conn.Resolver

	if resolvedConn.Err != nil {
		t.Close()
		return
	}

	client := resolvedConn.Client

	for {
		data, err := t.session.Recv()

		if err != nil {
			s.inactivateChan <- client
			return
		}

		s.incomingChan <- notify.IncomingData{client, strings.NewReader(data)}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) Accept() <-chan notify.IncomingConnection {
	return s.acceptChan
}

func (s *Server) Inactive() <-chan *pubsub.Client {
	return s.inactivateChan
}

func (s *Server) Incoming() <-chan notify.IncomingData {
	return s.incomingChan
}

func NewServer(prefix string, opts sockjs.Options) *Server {
	s := &Server{
		acceptChan:   make(chan notify.IncomingConnection),
		incomingChan: make(chan notify.IncomingData),
	}
	s.handler = sockjs.NewHandler(prefix, opts, s.serveSockJS)
	return s
}