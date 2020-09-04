package notify

import (
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify/message"
	"sync"
	"time"
)

type PubSub struct {
	clients 		map[string]*Client
	users 			map[string]map[string]struct{}
	channels 		map[string]*Channel
	subs 			map[string]map[string]struct{}
	config 			PubSubConfig
	inactiveClients map[string]time.Time
	app 			*App
	mu 				sync.RWMutex
}

// Connect adds a new client. Also it can restore client from garbage and send lost data.
func (pubsub *PubSub) Connect(info ClientInfo, transport Transport) (*Client, error) {
	if info.ID == NilId || transport == nil {
		return nil, errors.New("invalid client info")
	}
	pubsub.mu.Lock()
	client, ok := pubsub.clients[info.ID]
	if !ok {
		client = &Client{
			id:        info.ID,
			userId:    info.UserID,
			transport: transport,
			state:     ActiveClient,
			messageBuffer: message.NewBuffer(pubsub.config.ClientMessageBufferSize),
			app:       pubsub.app,
		}
		pubsub.clients[client.id] = client
		if info.UserID != NilId {
			if user, ok := pubsub.users[info.UserID]; ok {
				user[client.id] = struct{}{}
			} else {
				user = map[string]struct{}{}
				pubsub.users[info.UserID] = user
				user[info.ID] = struct{}{}
			}
		}
		pubsub.mu.Unlock()
		return client, nil
	}
	if client.userId != info.UserID {
		//pubsub.mu.Unlock()
		return nil, errors.New("client id does not match user id")
	}
	var err error = nil
	if ok, buffer := client.tryActivate(transport); ok {
		pubsub.mu.Unlock()
		delete(pubsub.inactiveClients, client.id)
		//todo: figure out what to do with not found messages
		fmt.Println("restoring...", buffer)
		messages, _ := pubsub.app.loadMessages(buffer)
		for i := 0; i < len(messages); i++ {
			message := messages[i]
			res := client.send(message)
			if res == invalidClientSend {
				i--
				pubsub.mu.RLock()
				client = pubsub.clients[client.id]
				pubsub.mu.RUnlock()
				if client == nil {
					break
				}
			}
		}
	} else {
		if client.tryInvalidate() {
			client = &Client{
				id:        info.ID,
				userId:    info.UserID,
				transport: transport,
				messageBuffer: message.NewBuffer(pubsub.config.ClientMessageBufferSize),
				state:     ActiveClient,
				app:       pubsub.app,
			}
			pubsub.clients[client.id] = client
			if info.UserID != NilId {
				if user, ok := pubsub.users[info.UserID]; ok {
					user[client.id] = struct{}{}
				} else {
					user = map[string]struct{}{}
					pubsub.users[info.UserID] = user
					user[info.ID] = struct{}{}
				}
			}
		} else {
			client = nil
			err = errors.New("could not connect")
		}
		pubsub.mu.Unlock()
	}
	return client, err

}

// Deletes client by id
func (pubsub *PubSub) Disconnect(id string) {
	if id == NilId {
		return
	}

	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	client, ok := pubsub.clients[id]
	if !ok {
		return
	}
	client.inactivate()
	pubsub.inactiveClients[client.id] = time.Now()
}

// DisconnectClient deletes client from pubsub by it's instance.
// It checks if instance is valid.
func (pubsub *PubSub) DisconnectClient(client *Client) {
	if client == nil {
		return
	}

	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	realClient := pubsub.clients[client.id]
	if realClient != client {
		return
	}
	client.inactivate()
	pubsub.inactiveClients[client.id] = time.Now()
}

// Join adds clients to specified channels
func (pubsub *PubSub) Join(opts JoinOptions) {
	clients := make([]string, len(opts.Clients))
	copy(clients, opts.Clients)

	pubsub.mu.RLock()
	if len(pubsub.users) > 0 {
		if len(opts.Users) > 0 {
			for _, id := range opts.Users {
				if user, ok := pubsub.users[id]; ok {
					for cid := range user {
						clients = append(clients, cid)
					}
				}
			}
		}
	}
	pubsub.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	for _, clientId := range clients {
		client, ok := pubsub.clients[clientId]
		if !ok {
			continue
		}
		subs, ok := pubsub.subs[client.id]
		if !ok || subs == nil {
			subs = map[string]struct{}{}
			pubsub.subs[client.id] = subs
		}

		for _, id := range opts.Channels {
			channel, ok := pubsub.channels[id]
			if !ok {
				channel = &Channel{
					members: nil,
					id:      id,
					state:   ActiveChannel,
				}
				pubsub.channels[id] = channel
			}
			channel.add(client)
			subs[channel.id] = struct{}{}
		}
	}
}

// Leave deletes clients from specified channels
func (pubsub *PubSub) Leave(opts LeaveOptions) {
	if !opts.All && len(opts.Channels) == 0 {
		return
	}
	clients := make([]string, 0, len(opts.Clients))
	copy(clients, opts.Clients)

	if len(opts.Users) > 0 {
		pubsub.mu.RLock()
		if len(pubsub.clients) > 0 {
			for _, id := range opts.Users {
				if user, ok := pubsub.users[id]; ok {
					for cid := range user {
						clients = append(clients, cid)
					}
				}
			}
		}
		pubsub.mu.RUnlock()
	}
	if len(clients) == 0 {
		return
	}

	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()

	if opts.All {
		for _, clientId := range clients {
			subs, ok := pubsub.subs[clientId]
			if !ok || len(subs) == 0{
				continue
			}
			pubsub.subs[clientId] = nil
			for sub := range subs {
				if channel, ok := pubsub.channels[sub]; ok {
					channel.del(clientId)
				}
			}
		}
	} else {
		for _, clientId := range clients {
			subs, ok := pubsub.subs[clientId]
			if !ok || len(subs) == 0{
				continue
			}
			for _, chID := range opts.Channels {
				if channel, ok := pubsub.channels[chID]; ok {
					delete(subs, chID)
					channel.del(clientId)
				}
			}
		}
	}
}

// Send sends Message to channels and clients
func (pubsub *PubSub) Send(opts SendOptions) {
	clients := make([]*Client, 0, len(opts.Clients)+len(opts.Users))
	channels := make([]*Channel, 0, len(opts.Channels))

	pubsub.mu.RLock()
	if len(pubsub.clients) > 0 {
		if len(opts.Clients) > 0 {
			for _, id := range opts.Clients {
				if client, ok := pubsub.clients[id]; ok {
					clients = append(clients, client)
				}
			}
		}
		if len(opts.Users) > 0 {
			for _, id := range opts.Users {
				if user, ok := pubsub.users[id]; ok {
					for cid := range user {
						if client, ok := pubsub.clients[cid]; ok {
							clients = append(clients, client)
						}
					}
				}
			}
		}
	}
	if len(opts.Channels) > 0 {
		for _, id := range opts.Channels {
			if channel, ok := pubsub.channels[id]; ok {
				channels = append(channels, channel)
			}
		}
	}
	pubsub.mu.RUnlock()
	for _, client := range clients {
		client.send(opts.Message)
	}
	for _, channel := range channels {
		channel.send(opts.Message)
	}
}

// Clean is supposed to delete clients that have been inactive for a specified time.
// After Clean deleted clients become invalid and cannot be used anymore.
// Returns list of ids of deleted clients
func (pubsub *PubSub) Clean() []string {
	pubsub.mu.RLock()
	inactiveIds := make([]string, len(pubsub.inactiveClients))
	now := time.Now()
	for id, t := range pubsub.inactiveClients {
		if now.Sub(t) > pubsub.config.ClientTTL {
			inactiveIds = append(inactiveIds, id)
		}
	}
	pubsub.mu.RUnlock()
	if len(inactiveIds) == 0 {
		return nil
	}
	clients := make([]string, 0, len(inactiveIds))
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	for _, id := range inactiveIds {
		if client, ok := pubsub.clients[id]; ok {
			if client.tryInvalidate() {
				clients = append(clients, id)
				delete(pubsub.clients, id)
				if client.userId != NilId {
					if user, ok := pubsub.users[client.userId]; ok {
						delete(user, id)
						if len(user) == 0 {
							delete(pubsub.users, client.userId)
						}
					}
				}
				if subs, ok := pubsub.subs[id]; ok {
					for sub, _ := range subs {
						if channel, ok := pubsub.channels[sub]; ok {
							channel.del(id)
						}
					}
					delete(pubsub.subs, id)
				}
			}
		}
	}
	return clients
}

func newPubsub(app *App, config PubSubConfig) *PubSub {
	if config.ClientTTL == 0 {
		config.ClientTTL = time.Minute * 5
	}
	pubsub := &PubSub{
		clients:  		 map[string]*Client{},
		users:   		 map[string]map[string]struct{}{},
		channels: 		 map[string]*Channel{},
		subs:     		 map[string]map[string]struct{}{},
		inactiveClients: map[string]time.Time{},
		config:   		 config,
		app:      		 app,
	}
	return pubsub
}