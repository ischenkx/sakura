package dnslb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RomanIschenko/notify"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

const PingEvent = "lb_ping_event"
const PongEvent = "lb_pong_event"

var lbLogger = logrus.WithField("source", "lblb")

type instance struct {
	ID, Address string
	addTime time.Time
}

type Pong struct {
	ID, BrokerID, IP string
}

type Config struct {
	Broker notify.Broker
	PingInterval time.Duration
	PongDeadline time.Duration
	PongsBufferSize int
}

func (cfg Config) validate() Config {
	if cfg.PingInterval <= 0 {
		cfg.PingInterval = time.Minute * 2
	}

	if cfg.PongsBufferSize <= 0 {
		cfg.PongsBufferSize = 128
	}

	if cfg.PongDeadline <= 0 {
		cfg.PongDeadline = time.Minute*3/2
	}

	return cfg
}

type loadBalancer struct {
	pingInterval time.Duration
	pongDeadline time.Duration
	broker notify.Broker
	roundRobin *roundRobin
	recentlyDeleted map[string]time.Time
	instances map[string]instance
	pongs chan Pong
	mu sync.Mutex
}

func (lb *loadBalancer) startPinging(ctx context.Context) {
	if lb.broker == nil {
		return
	}
	timeout := time.After(0)
	for {
		select {
		case <-timeout:
			pingTime := time.Now()
			id := uuid.New().String()
			lb.broker.Emit(notify.BrokerEvent{
				Data:  id,
				Event: PingEvent,
			})
			deadline := time.After(lb.pongDeadline)
			replies := map[string]string{}

			for {
				f := false
				select {
				case <-ctx.Done():
					return
				case pong := <-lb.pongs:
					if pong.ID == id {
						replies[pong.BrokerID] = pong.IP
						lb.mu.Lock()
						lb.add(instance{
							ID:        pong.BrokerID,
							Address:   pong.IP,
							addTime: pingTime,
						})
						lb.mu.Unlock()
					}
				case <-deadline:
					f = true
				}
				if f {
					break
				}
			}

			lb.mu.Lock()
			for brokerID := range lb.instances {
				if _, ok := replies[brokerID]; !ok {
					lb.del(brokerID, pingTime)
				}
			}
			lb.mu.Unlock()

			timeout = time.After(lb.pingInterval)
		case <-ctx.Done():
			return
		}
	}
}

func (lb *loadBalancer) handleBrokerEvents(ctx context.Context) {
	if lb.broker == nil {
		return
	}
	handle := lb.broker.Handle(func(e notify.BrokerEvent){
		switch e.Event {
		case notify.BrokerAppUpEvent:
			if ip, ok := e.Data.(string); ok {
				lb.mu.Lock()
				lb.add(instance{
					ID:        e.BrokerID,
					Address:   ip,
					addTime: time.Now(),
				})
				lb.mu.Unlock()
			}
		case notify.BrokerAppDownEvent:
			lb.mu.Lock()
			lb.del(e.BrokerID, time.Now())
			lb.mu.Unlock()
		case PongEvent:
			var pong Pong
			if str, ok := e.Data.(string); ok {
				if err := json.Unmarshal([]byte(str), &pong); err == nil {
					pong.BrokerID = e.BrokerID
					select {
					case lb.pongs <- pong:
					}
				}
			}
		}
	})
	defer handle.Close()
	<-ctx.Done()
}

func (lb *loadBalancer) next() (string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	id, err := lb.roundRobin.next()
	if err != nil {
		return "", err
	}
	if inst, ok := lb.instances[id]; ok {
		return inst.Address, nil
	}
	lbLogger.Debugf("failed to get address of %s", id)
	return "", errors.New("failed to get instance")
}

// adds instance, not thread-safe (requires lock)
func (lb *loadBalancer) add(inst instance) {
	lbLogger.Debugf("adding instance %s", inst.ID)
	if _, ok := lb.instances[inst.ID]; ok {
		return
	}

	if t, ok := lb.recentlyDeleted[inst.ID]; ok {
		if t.UnixNano() > inst.addTime.UnixNano() {
			return
		}
		delete(lb.recentlyDeleted, inst.ID)
	}

	lb.roundRobin.add(inst.ID)
	lb.instances[inst.ID] = inst
}

// deletes server instance by id, not thread-safe (requires lock)
func (lb *loadBalancer) del(id string, deletionTime time.Time) {
	inst, ok := lb.instances[id]
	if !ok {
		return
	}
	if inst.addTime.UnixNano() > deletionTime.UnixNano() {
		return
	}
	lbLogger.Debugf("deleting instance %s", id)
	delete(lb.instances, id)
	lb.recentlyDeleted[id] = deletionTime
	lb.roundRobin.del(id)
}

func (lb *loadBalancer) startCleaning(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 3)
	invalidationTime := int64(time.Minute * 2)
	for {
		now := time.Now().UnixNano()
		select {
		case <-ticker.C:
			lb.mu.Lock()
			for id, t := range lb.recentlyDeleted {
				if now - t.UnixNano() > invalidationTime {
					delete(lb.recentlyDeleted, id)
				}
			}
			lb.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

// Starts the lb server (non-blocking)
func (lb *loadBalancer) start(ctx context.Context) {
	go lb.startPinging(ctx)
	go lb.handleBrokerEvents(ctx)
	go lb.startCleaning(ctx)
}

func newLoadBalancer(cfg Config) *loadBalancer {
	cfg = cfg.validate()
	return &loadBalancer{
		broker:       cfg.Broker,
		pongs:        make(chan Pong, cfg.PongsBufferSize),
		pingInterval: cfg.PingInterval,
		pongDeadline: cfg.PongDeadline,
		instances:    map[string]instance{},
		recentlyDeleted: map[string]time.Time{},
		roundRobin:   newRoundRobin(),
	}
}

func handlePingEvent(b notify.Broker, ipProvider func() string) io.Closer {
	return b.Handle(func(e notify.BrokerEvent) {
		if e.Event == PingEvent {
			ip := ipProvider()
			if id, ok := e.Data.(string); ok {
				pong := Pong{
					ID:       id,
					IP:       ip,
				}

				data, err := json.Marshal(pong)
				if err != nil {
					return
				}
				b.Emit(notify.BrokerEvent{
					Data:  string(data),
					Event: PongEvent,
				})
			}
		}
	})
}
