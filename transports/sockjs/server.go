package sockjs

import (
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/igm/sockjs-go/sockjs"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var logger = logrus.WithField("source", "sockjs_server")

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
		logger.Errorf("failed to recv data, closing:", err)
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
		logger.Errorf("failed to resolve connection, closing:", resolvedConn.Err)
		t.Close()
		return
	}

	client := resolvedConn.Client

	for {
		data, err := t.session.Recv()

		if err != nil {
			//logger.Debugf("(sockjs)session.Recv failed, inactivating connection:", err)
			//fmt.Println("pushing client!!!")
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
		acceptChan:   make(chan notify.IncomingConnection, 1024),
		incomingChan: make(chan notify.IncomingData, 1024),
		inactivateChan: make(chan *pubsub.Client, 1024),
	}
	s.handler = sockjs.NewHandler(prefix, opts, s.serveSockJS)
	return s
}