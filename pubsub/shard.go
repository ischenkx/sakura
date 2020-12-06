package pubsub

import (
	"errors"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/google/uuid"
	"sync"
	"time"
)

// TODO
// Handle situations when publish fails (re-publish queue with some interval?)

type subscriptionState int

const (
	inactiveSub subscriptionState = iota
	activeSub
)

type subscription struct {
	lastTouch int64
	state 	  subscriptionState
}

const MinimalClientTTL = time.Second*5
const MinimalClientBufferSize = 50

type ClientConfig struct {
	TTL 		 	   time.Duration
	InvalidationTime time.Duration
	BufferSize 	   int
}

func (cfg *ClientConfig) validate() {
	if cfg.TTL <= MinimalClientTTL {
		cfg.TTL = MinimalClientTTL
	}

	if cfg.BufferSize < MinimalClientBufferSize {
		cfg.BufferSize = MinimalClientBufferSize
	}
}


type shard struct {
	id				string
	clients  		map[string]*Client
	users			map[string]map[*Client]struct{}
	topics 	 		map[string]topic
	// client id => unix nanoseconds
	// inactive clients are clients which are disconnected for a short time
	// (example: transport - websocket, internet issues cause client disconnections
	// for 5 seconds or something like that, it'd be inefficient to delete it and recreate again)
	inactiveClients map[string]int64
	// client id => unix nanoseconds
	// invalid client is a client that was recently deleted from inactive clients
	// due to timeout
	invalidClients  map[string]int64
	subs	 		map[string]map[string]*subscription
	userSubs 		map[string]map[string]*subscription
	nsRegistry 		*namespaceRegistry
	clientConfig	ClientConfig
	mu 				sync.RWMutex
}

func (s *shard) Publish(opts PublishOptions) {
	if len(opts.Clients) == 0 && len(opts.Topics) == 0 && len(opts.Users) == 0 {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, userID := range opts.Users {
		if user, ok := s.users[userID]; ok {
			for client := range user {
				client.publish(opts.Payload)
			}
		}
	}
	for _, clientID := range opts.Clients {
		if client, ok := s.clients[clientID]; ok {
			client.publish(opts.Payload)
		}
	}

	for _, topicID := range opts.Topics {
		if topic, ok := s.topics[topicID]; ok {
			for client := range topic.subscribers() {
				client.publish(opts.Payload)
			}
		}
	}
}

func (s *shard) Subscribe(opts SubscribeOptions) (res changelog.Log) {
	if (len(opts.Users) == 0 && len(opts.Clients) == 0) || len(opts.Topics) == 0 {
		return
	}
	res.Time = opts.Time
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, topicID := range opts.Topics {
		topic, ok := s.topics[topicID]
		if !ok {
			topic = s.nsRegistry.generate(topicID)
			s.topics[topicID] = topic
			res.TopicsUp = append(res.TopicsUp, topicID)
		}

		for _, clientID := range opts.Clients {
			if client, ok := s.clients[clientID]; ok {
				subs, ok := s.subs[clientID]
				if !ok {
					subs = map[string]*subscription{}
					s.subs[clientID] = subs
				}
				sub, ok := subs[topicID]
				if ok {
					if sub.lastTouch > opts.Time {
						continue
					}
				} else {
					sub = &subscription{}
					subs[topicID] = sub
				}
				sub.lastTouch = opts.Time
				sub.state = activeSub
				// todo: handle error
				topic.add(client)
			}
		}
		for _, userID := range opts.Users {
			if user, ok := s.users[userID]; ok {
				userSubs, ok := s.userSubs[userID]
				if !ok {
					userSubs = map[string]*subscription{}
					s.userSubs[userID] = userSubs
				}
				sub, ok := userSubs[topicID]
				if ok {
					if sub.lastTouch > opts.Time {
						continue
					}
				} else {
					sub = &subscription{}
				}

				sub.lastTouch = opts.Time
				sub.state = activeSub

				for client := range user {
					if clientSubs, ok := s.subs[client.ID()]; ok {
						if _, ok := clientSubs[topicID]; ok {
							continue
						}
					}
					topic.add(client)
				}
			}
		}
	}
	return
}

func (s *shard) Unsubscribe(opts UnsubscribeOptions) (res changelog.Log) {
	if (len(opts.Users) == 0 && len(opts.Clients) == 0) || (len(opts.Topics) == 0 && !opts.All) {
		return
	}
	res.Time = opts.Time
	s.mu.Lock()
	defer s.mu.Unlock()

	if opts.All {
		for _, clientID := range opts.Clients {
			subs, ok := s.subs[clientID]
			if !ok {
				continue
			}
			var client *Client = nil
			for topicID, sub := range subs {
				if sub.lastTouch > opts.Time {
					continue
				}

				if sub.state == inactiveSub {
					sub.lastTouch = opts.Time
					continue
				}

				if client == nil {
					client = s.clients[clientID]
				}
				if topic, ok := s.topics[topicID]; ok {
					if l, _ := topic.del(client); l == 0 {
						delete(s.topics, topicID)
						res.TopicsDown = append(res.TopicsDown, topicID)
					}
				}

				sub.state = inactiveSub
				sub.lastTouch = opts.Time
			}
		}

		for _, userID := range opts.Users {
			userSubs, ok := s.userSubs[userID]

			if !ok {
				continue
			}

			user := s.users[userID]

			for topicID, sub := range userSubs {
				if sub.lastTouch > opts.Time {
					continue
				}
				if sub.state == inactiveSub {
					sub.lastTouch = opts.Time
					continue
				}

				sub.state = inactiveSub
				sub.lastTouch = opts.Time

				topic, ok := s.topics[topicID]

				if !ok {
					continue
				}

				for client := range user {
					if clientSubs, ok := s.subs[client.ID()]; ok {
						if _, ok := clientSubs[topicID]; ok {
							continue
						}
					}
					if l, _ := topic.del(client); l == 0 {
						delete(s.topics, topicID)
						res.TopicsDown = append(res.TopicsDown, topicID)
						break
					}
				}
			}
		}
		return
	}

	topicLoop:
	for _, topicID := range opts.Topics {
		topic, ok := s.topics[topicID]
		if !ok {
			continue
		}
		for _, clientID := range opts.Clients {
			if subs, ok := s.subs[clientID]; ok {
				if sub, ok := subs[topicID]; ok {
					if sub.lastTouch > opts.Time {
						continue
					}
					if sub.state == inactiveSub {
						sub.lastTouch = opts.Time
						continue
					}
					sub.state = inactiveSub
					sub.lastTouch = opts.Time

					client := s.clients[clientID]

					if l, _ := topic.del(client); l == 0 {
						continue topicLoop
					}
				}
			}
		}

		for _, userID := range opts.Users {
			if userSubs, ok := s.userSubs[userID]; ok {
				if sub, ok := userSubs[topicID]; ok {
					if sub.lastTouch > opts.Time {
						continue
					}
					if sub.state == inactiveSub {
						sub.lastTouch = opts.Time
						continue
					}

					sub.lastTouch = opts.Time
					sub.state = inactiveSub

					user := s.users[userID]

					for client := range user {
						if clientSubs, ok := s.subs[client.ID()]; ok {
							if _, ok := clientSubs[topicID]; ok {
								continue
							}
						}
						if l, _ := topic.del(client); l == 0 {
							delete(s.topics, topicID)
							res.TopicsDown = append(res.TopicsDown, topicID)
							continue topicLoop
						}
					}
				}
			}
		}
	}
	return
}

func (s *shard) Connect(opts ConnectOptions) (*Client, changelog.Log, error) {
	var res changelog.Log
	if opts.Transport == nil {
		return nil, res, errors.New("connect failed: transport is nil")
	}
	if opts.ID == "" {
		return nil, res, errors.New("connect failed: client id is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, invalid := s.invalidClients[opts.ID]; invalid {
		return nil, res, errors.New("client has been invalidated, try to reconnect")
	}
	c, clientExisted := s.clients[opts.ID]

	if !clientExisted {
		var err error
		c, err = newClient(opts.ID, s.clientConfig.BufferSize)
		if err != nil {
			return nil, changelog.Log{}, err
		}
	} else if c.UserID() != GetUserID(opts.ID) {
		return nil, changelog.Log{}, errors.New("invalid user id")
	}

	if err := c.tryActivate(opts.Transport); err == nil {

		s.clients[c.ID()] = c
		delete(s.inactiveClients, c.ID())

		userID := GetUserID(opts.ID)
		if userID != "" {
			user, ok := s.users[userID]
			if ok {
				user[c] = struct{}{}
			} else {
				user = map[*Client]struct{}{}
				user[c] = struct{}{}
				s.users[userID] = user
				res.UsersUp = append(res.UsersUp, userID)
			}
			if !clientExisted {
				if userSubs, ok := s.userSubs[userID]; ok {
					for topicID, sub := range userSubs {
						if sub.state == activeSub {
							s.topics[topicID].add(c)
						}
					}
				}
			}
		}

		res.ClientsUp = append(res.ClientsUp, c.ID())

		return c, res, nil
	} else {
		return nil, res, err
	}
}

func (s *shard) InactivateClient(client *Client) (res changelog.Log) {
	if client == nil {
		return
	}
	res.Time = time.Now().UnixNano()

	s.mu.Lock()
	defer s.mu.Unlock()
	if realClient, exists := s.clients[client.ID()]; exists {
		if realClient != client {
			return
		}
		if err := client.tryInactivate(); err != nil {
			return
		}
		res.ClientsDown = append(res.ClientsDown, client.ID())
		s.inactiveClients[client.ID()] = time.Now().UnixNano()
	}
	return
}

func (s *shard) Disconnect(opts DisconnectOptions) (res changelog.Log) {
	if len(opts.Users) == 0 && len(opts.Clients) == 0 && !opts.All {
		return
	}
	res.Time = opts.Time
	if opts.All {
		for _, client := range s.clients {
			client.forceInvalidate()
		}
		res.TopicsDown = make([]string, 0, len(s.topics))

		for topicID := range s.topics {
			res.TopicsDown = append(res.TopicsDown, topicID)
		}
		s.clients = map[string]*Client{}
		s.users = map[string]map[*Client]struct{}{}
		s.topics = map[string]topic{}
		s.invalidClients = map[string]int64{}
		s.inactiveClients = map[string]int64{}
		s.userSubs = map[string]map[string]*subscription{}
		s.subs = map[string]map[string]*subscription{}
		return
	}

	for _, userID := range opts.Users {
		if user, ok := s.users[userID]; ok {
			delete(s.users, userID)
			userSubs := s.userSubs[userID]
			delete(s.userSubs, userID)
			for client := range user {
				client.forceInvalidate()
				if userSubs != nil {
					for topicID, sub := range userSubs {
						if sub.state == activeSub {
							topic := s.topics[topicID]
							if l, _ := topic.del(client); l == 0 {
								res.TopicsDown = append(res.TopicsDown, topicID)
								delete(s.topics, topicID)
							}
						}
					}
				}
				if subs, ok := s.subs[client.ID()]; ok {
					delete(s.subs, client.ID())
					for topicID, sub := range subs {
						if sub.state == activeSub {
							topic := s.topics[topicID]
							if l, _ := topic.del(client); l == 0 {
								res.TopicsDown = append(res.TopicsDown, topicID)
								delete(s.topics, topicID)
							}
						}
					}
				}

				s.invalidClients[client.ID()] = opts.Time
				delete(s.inactiveClients, client.ID())
			}
		}
	}

	for _, clientID := range opts.Clients {
		if client, ok := s.clients[clientID]; ok {
			client.forceInvalidate()
			delete(s.clients, client.ID())
			if subs, ok := s.subs[clientID]; ok {
				delete(s.subs, clientID)
				for topicID, sub := range subs {
					if sub.state == activeSub {
						topic := s.topics[topicID]
						if l, _ := topic.del(client); l == 0 {
							res.TopicsDown = append(res.TopicsDown, topicID)
							delete(s.topics, topicID)
						}
					}
				}
			}
			if userID := client.UserID(); userID != "" {
				if user, ok := s.users[userID]; ok {
					delete(user, client)
					if userSubs, ok := s.userSubs[userID]; ok {
						for topicID, sub := range userSubs {
							if sub.state == activeSub {
								topic := s.topics[topicID]
								if l, _ := topic.del(client); l == 0 {
									res.TopicsDown = append(res.TopicsDown, topicID)
									delete(s.topics, topicID)
								}
							}
						}
					}
					if len(user) == 0 {
						delete(s.users, userID)
						delete(s.userSubs, userID)
					}
				}
			}
		}
	}
	return
}

func (s *shard) Clean() (res changelog.Log) {
	var inactiveClients []string
	var invalidClients []string
	now := time.Now().UnixNano()
	invalidationTime := int64(s.clientConfig.InvalidationTime)
	clientTTl := int64(s.clientConfig.TTL)
	res.Time = now
	s.mu.RLock()
	for id, t := range s.inactiveClients {
		if now - t >= clientTTl {
			inactiveClients = append(inactiveClients, id)
		}
	}
	for id, t := range s.invalidClients {
		if now - t >= invalidationTime {
			invalidClients = append(invalidClients, id)
		}
	}
	s.mu.RUnlock()

	if len(inactiveClients) == 0 && len(invalidClients) == 0 {
		return
	}
	s.mu.Lock()
	for _, id := range invalidClients {
		res.ClientsDown = append(res.ClientsDown, id)
		delete(s.invalidClients, id)
	}
	s.mu.Unlock()
	for _, id := range inactiveClients {
		s.mu.Lock()
		if client, ok := s.clients[id]; ok {
			if err := client.tryInvalidate(); err == nil {
				delete(s.clients, id)
				res.ClientsDown = append(res.ClientsDown, id)
				if subs, ok := s.subs[id]; ok {
					delete(s.subs, id)
					for topicID, sub := range subs {
						if sub.state == activeSub {
							topic := s.topics[topicID]
							if l, _ := topic.del(client); l == 0 {
								res.TopicsDown = append(res.TopicsDown, topicID)
								delete(s.topics, topicID)
							}
						}
					}
				}
				if userID := client.UserID(); userID != "" {
					if user, ok := s.users[userID]; ok {
						delete(user, client)
						if userSubs, ok := s.userSubs[userID]; ok {
							for topicID, sub := range userSubs {
								if sub.state == activeSub {
									topic := s.topics[topicID]
									if l, _ := topic.del(client); l == 0 {
										res.TopicsDown = append(res.TopicsDown, topicID)
										delete(s.topics, topicID)
									}
								}
							}
						}
						if len(user) == 0 {
							res.UsersDown = append(res.UsersDown, userID)
							delete(s.users, userID)
							delete(s.userSubs, userID)
						}
					}
				}
			}
		}
		delete(s.inactiveClients, id)
		s.mu.Unlock()
	}
	return
}

func (s *shard) Clients() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	clients := make([]string, 0, len(s.clients))
	for id := range s.clients {
		clients = append(clients, id)
	}
	return clients
}

func (s *shard) Users() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]string, 0, len(s.users))
	for id := range s.users {
		users = append(users, id)
	}
	return users
}

func (s *shard) Metrics() (m Metrics) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m.Clients = len(s.clients)
	m.Users = len(s.users)
	return
}

func newShard(nsRegistry *namespaceRegistry, clientCfg ClientConfig) *shard {
	clientCfg.validate()
	return &shard{
		id: uuid.New().String(),
		clients: 		 map[string]*Client{},
		topics: 		 map[string]topic{},
		inactiveClients: map[string]int64{},
		invalidClients:  map[string]int64{},
		subs: 			 map[string]map[string]*subscription{},
		users: 			 map[string]map[*Client]struct{}{},
		userSubs: 		 map[string]map[string]*subscription{},
		nsRegistry: 	 nsRegistry,
		clientConfig:    clientCfg,
	}
}