package websockets_notify

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"net/http"
	"notify"
	"sync"
	"sync/atomic"
)

const (
	ClosedWS int32 = iota + 1
	OpenWS
)

type WSServer struct {
	upgrader ws.HTTPUpgrader
	notifyServer notify.Server
	auth Authenticator
	mu sync.RWMutex
}

type WSTransport struct {
	conn net.Conn
	state int32
}

func (t *WSTransport) Send(mes notify.Message) {
	if err := wsutil.WriteServerBinary(t.conn, mes.Data); err != nil {
		t.Close()
	}
}

func (t *WSTransport) Close() {
	if atomic.CompareAndSwapInt32(&t.state, OpenWS, ClosedWS) {
		t.conn.Close()
	}
}

func (handler *WSServer) serveConn(client *notify.Client, t *WSTransport) {
	defer t.Close()
	state := ws.StateServerSide
	ch := wsutil.ControlFrameHandler(t.conn, state)
	r := &wsutil.Reader{
		Source:         t.conn,
		State:          state,
		CheckUTF8:      true,
		OnIntermediate: ch,
	}
	for {
		if atomic.LoadInt32(&t.state) == ClosedWS {
			handler.notifyServer.DisconnectClient(client)
			return
		}
		h, err := r.NextFrame()
		if err != nil {
			handler.notifyServer.DisconnectClient(client)
			return
		}
		if h.OpCode.IsControl() {
			if err = ch(h, r); err != nil {
				return
			}
			continue
		}
		if handler.notifyServer.Handle(client, r) != nil {
			return
		}
		if err := r.Discard(); err != nil {
			fmt.Println(err)
		}
	}
}

func (handler *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler.notifyServer == nil {
		log.Println("servable is nil, server won't work without any servable")
		return
	}
	clientData, err := handler.auth.Authenticate(r)
	if err != nil {
		return
	}
	if conn, _, _, err := handler.upgrader.Upgrade(r, w); err == nil {
		t := &WSTransport{conn, OpenWS}
		info := notify.ClientInfo{
			ID:        clientData.ID,
			UserID:    clientData.UserID,
			Transport: t,
		}
		if client, err := handler.notifyServer.Connect(clientData.App, info); err == nil {
			go handler.serveConn(client, t)
		} else {
			t.Close()
		}
	} else {
		if conn != nil {
			if err := conn.Close(); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func New(servable notify.Server, auth Authenticator) *WSServer {
	if auth == nil {
		auth = DefaultAuthenticator{}
	}
	return &WSServer{
		upgrader: ws.HTTPUpgrader{},
		auth:     auth,
		notifyServer: servable,
	}
}