package notify

import (
	"fmt"
	"sync"
)

type ClientInfo struct {
	ID string
	UserID string
	Transport Transport
}

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
	messageBuffer MessageBuffer
	mu sync.Mutex
	state ClientState
	app *App
}

type clientSendResult int

const (
	successfulSend clientSendResult = iota + 1
	invalidClientSend = iota + 1
	inactiveClientSend = iota + 1

)

func (client *Client) ID() string {
	return client.id
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

func (client *Client) IsActive() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.state == ActiveClient
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

func (client *Client) send(mes Message) clientSendResult {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.state == InvalidClient {
		return invalidClientSend
	}
	if client.state == InactiveClient {
		client.messageBuffer.Push(mes.ID)
		if len(client.messageBuffer.buffer) % 1000 == 0 {
			fmt.Println(len(client.messageBuffer.buffer))
		}
		return inactiveClientSend
	}
	client.transport.Send(mes)
	return successfulSend
}

func (client *Client) Send(data []byte) {
	client.app.Send(MessageOptions{
		Clients:  []string{client.id},
		Data:     data,
	})
}

func (client *Client) Join(channels []string) {
	if len(channels) == 0 {
		return
	}
	client.app.Join(JoinOptions{
		Clients:  []string{client.id},
		Channels: channels,
	})
}

func (client *Client) Leave(channels []string, all bool) {
	if len(channels) == 0 && !all {
		return
	}
	if all {
		client.app.Leave(LeaveOptions{
			Clients: []string{client.id},
			All: 	 true,
		})
		return
	}
	client.app.Leave(LeaveOptions{
		Clients:  []string{client.id},
		Channels: channels,
	})
}
