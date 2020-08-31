package notify

import (
	"github.com/RomanIschenko/notify/message"
	"github.com/RomanIschenko/notify/options"
	"sync"
)

type ClientState int

const (
	ActiveClient ClientState = iota + 1
	InactiveClient
	InvalidClient
)

type Client struct {
	id, userId string
	transport Transport

	//it is used to save messages while client is closed
	messageBuffer message.Buffer
	mu            sync.Mutex
	state         ClientState
	app           *App
	data          sync.Map
}

type clientSendResult int

const (
	successfulSend clientSendResult = iota + 1
	invalidClientSend
	inactiveClientSend
)

func (client *Client) UserID() string {
	return client.userId
}

func (client *Client) ID() string {
	return client.id
}

func (client *Client) IsActive() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.state == ActiveClient
}

func (client *Client) inactivate() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.state == ActiveClient {
		client.transport.Close()
		client.transport = nil
		client.state = InactiveClient
	}
}

func (client *Client) tryInvalidate() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.state == InactiveClient {
		client.state = InvalidClient
		return true
	}
	return client.state == InvalidClient
}

func (client *Client) tryActivate(t Transport) (bool, []string) {
	if t == nil {
		return false, nil
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.state != InactiveClient {
		return false, nil
	}
	if client.transport != nil {
		client.transport.Close()
	}
	client.transport = t
	client.state = ActiveClient
	buffer, _ := client.messageBuffer.Reset()
	return true, buffer
}

func (client *Client) send(mes message.Message) clientSendResult {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.state == InvalidClient {
		return invalidClientSend
	}
	if client.state == InactiveClient {
		client.messageBuffer.Push(mes.ID)
		return inactiveClientSend
	}
	client.transport.Send(mes)
	return successfulSend
}

func (client *Client) Send(data []byte) {
	client.app.SendMessage(options.MessageSend{
		Clients:  []string{client.id},
		Data:     data,
	})
}

func (client *Client) Join(channels []string) {
	if len(channels) == 0 {
		return
	}
	client.app.Join(options.Join{
		Clients:  []string{client.id},
		Channels: channels,
	})
}

func (client *Client) Leave(channels []string, all bool) {
	if len(channels) == 0 && !all {
		return
	}
	client.app.Leave(options.Leave{
		Clients:      []string{client.id},
		Channels:     channels,
		All:          all,
	})
}

func (client *Client) Get(key string)(interface{}, bool) {
	return client.data.Load(key)
}

func (client *Client) Set(key string, value interface{}) {
	client.data.Store(key, value)
}

func (client *Client) Del(key string) {
	client.data.Delete(key)
}