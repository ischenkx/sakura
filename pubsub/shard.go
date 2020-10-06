package pubsub

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

// TODO
// Handle situations when publish fails (re-publish queue with some interval?)

const MinimalClientTTL = int64(time.Second*5)
const MinimalClientBufferSize = 50

type result struct {
	topicsUp, topicsDown []string
}

type ShardConfig struct {
	ClientTTL 		 	   int64
	ClientInvalidationTime int64
	ClientBufferSize 	   int
}

func (cfg ShardConfig) validate() ShardConfig {
	if cfg.ClientTTL <= MinimalClientTTL {
		cfg.ClientTTL = MinimalClientTTL
	}

	if cfg.ClientBufferSize < MinimalClientBufferSize {
		cfg.ClientBufferSize = MinimalClientBufferSize
	}

	return cfg
}

type subscriptionState int

const (
	inactiveSub subscriptionState = iota
	activeSub
)

type subscription struct {
	lastTouch int64
	state 	  subscriptionState
}

type shard struct {
	id				string
	clients  		map[string]*Client
	users			map[string]map[*Client]struct{}
	topics 	 		map[string]topic
	inactiveClients map[string]int64
	invalidClients  map[string]int64
	subs	 		map[string]map[string]*subscription
	userSubs 		map[string]map[string]*subscription
	nsRegistry 		*namespaceRegistry
	config			ShardConfig
	queue			pubQueue
	mu 				sync.RWMutex
}

func (s *shard) Publish(opts PublishOptions) (res result) {
	if len(opts.Clients) == 0 && len(opts.Topics) == 0 && len(opts.Users) == 0 {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, userID := range opts.Users {
		if user, ok := s.users[userID]; ok {
			for client := range user {
				//client.Publish(opts.Publication)
				s.queue.Enqueue(client, opts.Payload)
			}
		}
	}
	for _, clientID := range opts.Clients {
		if client, ok := s.clients[clientID]; ok {
			s.queue.Enqueue(client, opts.Payload)
		}
	}

	for _, topicID := range opts.Topics {
		if topic, ok := s.topics[topicID]; ok {
			for client := range topic.subscribers() {
				s.queue.Enqueue(client, opts.Payload)
			}
		}
	}
	return
}

func (s *shard) Subscribe(opts SubscribeOptions) (res result) {
	if (len(opts.Users) == 0 && len(opts.Clients) == 0) || len(opts.Topics) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, topicID := range opts.Topics {
		topic, ok := s.topics[topicID]
		if !ok {
			topic = s.nsRegistry.generate(topicID)
			s.topics[topicID] = topic
			res.topicsUp = append(res.topicsUp, topicID)
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
					if clientSubs, ok := s.subs[string(client.ID())]; ok {
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

func (s *shard) Unsubscribe(opts UnsubscribeOptions) (res result) {
	if (len(opts.Users) == 0 && len(opts.Clients) == 0) || (len(opts.Topics) == 0 && !opts.All) {
		return
	}
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
						res.topicsDown = append(res.topicsDown, topicID)
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
					if clientSubs, ok := s.subs[string(client.ID())]; ok {
						if _, ok := clientSubs[topicID]; ok {
							continue
						}
					}
					if l, _ := topic.del(client); l == 0 {
						delete(s.topics, topicID)
						res.topicsDown = append(res.topicsDown, topicID)
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
						if clientSubs, ok := s.subs[string(client.ID())]; ok {
							if _, ok := clientSubs[topicID]; ok {
								continue
							}
						}
						if l, _ := topic.del(client); l == 0 {
							delete(s.topics, topicID)
							res.topicsDown = append(res.topicsDown, topicID)
							continue topicLoop
						}
					}
				}
			}
		}
	}
	return
}

func (s *shard) Connect(opts ConnectOptions) (*Client, error) {
	if opts.Transport == nil {
		return nil, errors.New("connect failed: transport is nil")
	}
	if opts.ID == "" {
		return nil, errors.New("connect failed: client id is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, invalid := s.invalidClients[string(opts.ID)]; invalid {
		return nil, errors.New("client has been invalidated, try to reconnect")
	}

	client, clientExisted := s.clients[string(opts.ID)]

	if !clientExisted {
		var err error
		client, err = newClient(opts.ID, s.config.ClientBufferSize)
		if err != nil {
			return nil, err
		}
	} else if client.ID().User() != opts.ID.User() {
		return nil, errors.New("invalid user id")
	}

	if err := client.tryActivate(opts.Transport); err == nil {
		s.clients[string(client.ID())] = client
		delete(s.inactiveClients, string(client.ID()))

		if opts.ID.User() != "" {
			userID := opts.ID.User()
			user, ok := s.users[userID]
			if ok {
				user[client] = struct{}{}
			} else {
				user = map[*Client]struct{}{}
				user[client] = struct{}{}
				s.users[userID] = user
			}
			if !clientExisted {
				if userSubs, ok := s.userSubs[userID]; ok {
					for topicID, sub := range userSubs {
						if sub.state == activeSub {
							s.topics[topicID].add(client)
						}
					}
				}
			}
		}
		return client, nil
	} else {
		return nil, err
	}
}

func (s *shard) InactivateClient(client *Client) {
	if client == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if realClient, exists := s.clients[string(client.ID())]; exists {
		if realClient != client {
			return
		}
		if err := client.tryInactivate(); err != nil {
			return
		}
		s.inactiveClients[string(client.ID())] = time.Now().UnixNano()
	}
}

func (s *shard) Disconnect(opts DisconnectOptions) (res result) {
	if len(opts.Users) == 0 && len(opts.Clients) == 0 && !opts.All {
		return
	}
	if opts.All {
		for _, client := range s.clients {
			client.forceInvalidate()
		}
		res.topicsDown = make([]string, 0, len(s.topics))

		for topicID := range s.topics {
			res.topicsDown = append(res.topicsDown, topicID)
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
								res.topicsDown = append(res.topicsDown, topicID)
								delete(s.topics, topicID)
							}
						}
					}
				}
				if subs, ok := s.subs[string(client.ID())]; ok {
					delete(s.subs, string(client.ID()))
					for topicID, sub := range subs {
						if sub.state == activeSub {
							topic := s.topics[topicID]
							if l, _ := topic.del(client); l == 0 {
								res.topicsDown = append(res.topicsDown, topicID)
								delete(s.topics, topicID)
							}
						}
					}
				}

				s.invalidClients[string(client.ID())] = opts.Time
				delete(s.inactiveClients, string(client.ID()))
			}
		}
	}

	for _, clientID := range opts.Clients {
		if client, ok := s.clients[clientID]; ok {
			client.forceInvalidate()
			delete(s.clients, string(client.ID()))
			if subs, ok := s.subs[clientID]; ok {
				delete(s.subs, clientID)
				for topicID, sub := range subs {
					if sub.state == activeSub {
						topic := s.topics[topicID]
						if l, _ := topic.del(client); l == 0 {
							res.topicsDown = append(res.topicsDown, topicID)
							delete(s.topics, topicID)
						}
					}
				}
			}
			if userID := client.ID().User(); userID != "" {
				if user, ok := s.users[userID]; ok {
					delete(user, client)
					if userSubs, ok := s.userSubs[userID]; ok {
						for topicID, sub := range userSubs {
							if sub.state == activeSub {
								topic := s.topics[topicID]
								if l, _ := topic.del(client); l == 0 {
									res.topicsDown = append(res.topicsDown, topicID)
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

func (s *shard) Clean() (res result) {
	inactiveClients := []string{}
	invalidClients := []string{}
	now := time.Now().UnixNano()
	invalidationTime := s.config.ClientInvalidationTime
	clientTTl := s.config.ClientTTL
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
		delete(s.invalidClients, id)
	}
	s.mu.Unlock()
	for _, id := range inactiveClients {
		s.mu.Lock()
		if client, ok := s.clients[id]; ok {
			if err := client.tryInvalidate(); err == nil {
				delete(s.clients, id)
				if subs, ok := s.subs[id]; ok {
					delete(s.subs, id)
					for topicID, sub := range subs {
						if sub.state == activeSub {
							topic := s.topics[topicID]
							if l, _ := topic.del(client); l == 0 {
								res.topicsDown = append(res.topicsDown, topicID)
								delete(s.topics, topicID)
							}
						}
					}
				}
				if userID := client.ID().User(); userID != "" {
					if user, ok := s.users[userID]; ok {
						delete(user, client)
						if userSubs, ok := s.userSubs[userID]; ok {
							for topicID, sub := range userSubs {
								if sub.state == activeSub {
									topic := s.topics[topicID]
									if l, _ := topic.del(client); l == 0 {
										res.topicsDown = append(res.topicsDown, topicID)
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
		delete(s.inactiveClients, id)
		s.mu.Unlock()
	}
	return
}

func newShard(queue pubQueue, nsRegistry *namespaceRegistry, config ShardConfig) *shard {
	config = config.validate()
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
		config:          config,
		queue:			 queue,
	}
}