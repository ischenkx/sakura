package balancer

import (
	"context"
	"encoding/json"
	"github.com/RomanIschenko/notify/cluster/broker"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const PingEvent = "lb_ping_event"
const PongEvent = "lb_pong_event"
const InstanceUpEvent = "instance_up"
const InstanceDownEvent = "instance_down"
const GlobalChannel = "global_lb_dns"

var lbLogger = logrus.WithField("source", "load balancer")

type Pong struct {
	ID, BrokerID, IP string
}

type Config struct {
	Policy          Policy
	PingInterval    time.Duration
	PongDeadline    time.Duration
	PongsBufferSize int
}

func (cfg *Config) validate() {
	if cfg.Policy == nil {
		panic("invalid config: no load balancer_pb policy provided")
	}

	if cfg.PingInterval <= 0 {
		cfg.PingInterval = time.Second * 5
	}

	if cfg.PongsBufferSize <= 0 {
		cfg.PongsBufferSize = 128
	}

	if cfg.PongDeadline <= 0 {
		cfg.PongDeadline = time.Second * 15
	}
}

type Balancer struct {
	pingInterval    time.Duration
	pongDeadline    time.Duration
	pongs           chan Pong
	policy 			Policy
	recentlyDeleted map[string]time.Time
	instances       []Instance
	mu              sync.RWMutex
}

func (balancer *Balancer) startPinging(ctx context.Context, b broker.Broker) {
	if balancer == nil {
		return
	}
	timeout := time.After(0)
	for {
		select {
		case <-timeout:
			pingTime := time.Now()
			id := uuid.New().String()
			b.Publish([]string{GlobalChannel}, broker.NewEvent(PingEvent, []byte(id)))
			deadline := time.After(balancer.pongDeadline)
			replies := map[string]string{}
			for {
				f := false
				select {
				case <-ctx.Done():
					return
				case pong := <-balancer.pongs:
					if pong.ID == id {
						replies[pong.BrokerID] = pong.IP
						balancer.mu.Lock()
						balancer.add(Instance{
							ID:      pong.BrokerID,
							Address: pong.IP,
							time:    pingTime,
						})
						balancer.mu.Unlock()
					}
				case <-deadline:
					f = true
				}
				if f {
					break
				}
			}
			balancer.mu.Lock()
			for _, inst := range balancer.instances {
				brokerID := inst.ID
				if _, ok := replies[brokerID]; !ok {
					balancer.delete(brokerID, pingTime)
				}
			}
			balancer.mu.Unlock()

			timeout = time.After(balancer.pingInterval)
		case <-ctx.Done():
			return
		}
	}
}

func (balancer *Balancer) handleBrokerEvents(brk broker.Broker) {
	if brk == nil {
		return
	}

	brk.Subscribe([]string{GlobalChannel}, time.Now().UnixNano())

	brk.Handle(func(e broker.Event){
		if e.Name == InstanceDownEvent {
			balancer.mu.Lock()
			defer balancer.mu.Unlock()
			balancer.delete(e.BrokerID, time.Unix(0, e.Time))
		}
		if e.Name == InstanceUpEvent {
			balancer.mu.Lock()
			defer balancer.mu.Unlock()
			balancer.add(Instance{
				ID:      e.BrokerID,
				Address: string(e.Data),
				time:    time.Unix(0, e.Time),
			})
		}
		if e.Name == PongEvent {
			var pong Pong
			if err := json.Unmarshal(e.Data, &pong); err == nil {
				pong.BrokerID = e.BrokerID
				select {
				case balancer.pongs <- pong:
				}
			}
		}
	})
}

// adds Instance, not thread-safe (requires lock)
func (balancer *Balancer) add(inst Instance) {
	lbLogger.Debugf("adding Instance %s", inst.ID)
	for _, i := range balancer.instances {
		if i.ID == inst.ID {
			return
		}
	}
	if t, ok := balancer.recentlyDeleted[inst.ID]; ok {
		if t.UnixNano() > inst.time.UnixNano() {
			return
		}
		delete(balancer.recentlyDeleted, inst.ID)
	}
	balancer.instances = append(balancer.instances, inst)
}

func (balancer *Balancer) GetAddr() (string, error) {
	balancer.mu.RLock()
	defer balancer.mu.RUnlock()
	inst, err := balancer.policy.Next(addrList(balancer.instances))
	if err != nil {
		return "", err
	}
	return inst.Address, err
}

// deletes server Instance by id, not thread-safe (requires lock)
func (balancer *Balancer) delete(id string, deletionTime time.Time) {
	index := -1
	for i, inst := range balancer.instances {
		if inst.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	inst := balancer.instances[index]
	if inst.time.UnixNano() > deletionTime.UnixNano() {
		return
	}
	lbLogger.Debugf("deleting Instance %s", id)
	balancer.instances[index] = balancer.instances[len(balancer.instances)-1]
	balancer.instances = balancer.instances[:len(balancer.instances)-1]
	balancer.recentlyDeleted[id] = deletionTime
}

func (balancer *Balancer) startCleaning(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 3)
	invalidationTime := int64(time.Minute * 2)
	for {
		now := time.Now().UnixNano()
		select {
		case <-ticker.C:
			balancer.mu.Lock()
			for id, t := range balancer.recentlyDeleted {
				if now - t.UnixNano() > invalidationTime {
					delete(balancer.recentlyDeleted, id)
				}
			}
			balancer.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

// Starts the lb server (non-blocking)
func (balancer *Balancer) Start(ctx context.Context, broker broker.Broker) {
	balancer.handleBrokerEvents(broker)
	go balancer.startPinging(ctx, broker)
	go balancer.startCleaning(ctx)
}

func New(cfg Config) *Balancer {
	cfg.validate()
	return &Balancer{
		pongs:           make(chan Pong, cfg.PongsBufferSize),
		pingInterval:    cfg.PingInterval,
		pongDeadline:    cfg.PongDeadline,
		instances:       []Instance{},
		recentlyDeleted: map[string]time.Time{},
		policy: cfg.Policy,
	}
}