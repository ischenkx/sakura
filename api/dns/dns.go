package dns

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RomanIschenko/notify"
	"github.com/google/uuid"
	"io"
	"sync"
	"time"
)

const PingEvent = "dns_ping_event"
const PongEvent = "dns_pong_event"

type Pong struct {
	ID, BrokerID, IP string
}

type Config struct {
	Broker notify.Broker
	PingInterval time.Duration
	PongsBufferSize int
}

type DNS struct {
	pingInterval time.Duration
	broker notify.Broker
	roundRobin *roundRobin
	instances map[string]string
	pongs chan Pong
	mu sync.Mutex
}

func (dns *DNS) Next() (string, error) {
	dns.mu.Lock()
	defer dns.mu.Unlock()
	id, err := dns.roundRobin.next()
	if err != nil {
		return "", err
	}
	if ip, ok := dns.instances[id]; ok {
		return ip, nil
	}
	return "", errors.New("failed to get instance")
}

func (dns *DNS) Add(id, ip string) {
	dns.mu.Lock()
	defer dns.mu.Unlock()
	dns.roundRobin.add(id)
	dns.instances[id] = ip
}

func (dns *DNS) Del(id string) {
	dns.mu.Lock()
	defer dns.mu.Unlock()
	delete(dns.instances, id)
	dns.roundRobin.del(id)
}

func (dns *DNS) startPinging(ctx context.Context) {
	if dns.broker == nil {
		return
	}
	timeout := time.After(0)
	for {
		select {
		case <-timeout:
			id := uuid.New().String()
			dns.broker.Emit(notify.BrokerEvent{
				Data:     id,
				Event: PingEvent,
			})
			deadline := time.After(time.Minute * 3)
			replies := map[string]string{}
			pingLoop:
			for {
				select {
				case <-ctx.Done():
					return
				case pong := <-dns.pongs:
					if pong.ID == id {
						replies[pong.BrokerID] = pong.IP
						dns.Add(pong.BrokerID, pong.IP)
					}
				case <-deadline:
					break pingLoop
				}
			}

			dns.mu.Lock()
			for brokerID := range dns.instances {
				if _, ok := replies[brokerID]; !ok {
					delete(dns.instances, brokerID)
				}
			}
			dns.mu.Unlock()

			timeout = time.After(dns.pingInterval)
		case <-ctx.Done():
			return
		}
	}
}

func (dns *DNS) handleBrokerEvents(ctx context.Context) {
	if dns.broker == nil {
		return
	}

	handle := dns.broker.Handle(func(e notify.BrokerEvent){
		switch e.Event {
		case notify.BrokerAppUpEvent:
			if ip, ok := e.Data.(string); ok {
				dns.Add(e.BrokerID, ip)
			}
		case notify.BrokerAppDownEvent:
			dns.Del(e.BrokerID)
		case PongEvent:
			var pong Pong
			if str, ok := e.Data.(string); ok {
				if err := json.Unmarshal([]byte(str), &pong); err == nil {
					pong.BrokerID = e.BrokerID
					select {
					case dns.pongs <- pong:
					}
				}
			}
		}
	})

	defer handle.Close()

	<-ctx.Done()
}

func (dns *DNS) Start(ctx context.Context) {
	go dns.startPinging(ctx)
	go dns.handleBrokerEvents(ctx)
	<-ctx.Done()
}

func New(cfg Config) *DNS {
	return &DNS{
		broker: cfg.Broker,
		pongs: make(chan Pong, cfg.PongsBufferSize),
		pingInterval: cfg.PingInterval,
		instances: map[string]string{},
		roundRobin: newRoundRobin(),
	}
}

func WithIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, "appUpArg", ip)
}
func HandleDNSPingEvent(b notify.Broker, ipProvider func() string) io.Closer {
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
