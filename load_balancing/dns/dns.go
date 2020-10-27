package dnslb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

const PingEvent = "dns_ping_event"
const PongEvent = "dns_pong_event"

var dnsLogger = logrus.WithField("source", "dnslb")

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

type DNS struct {
	pingInterval time.Duration
	pongDeadline time.Duration
	broker notify.Broker
	roundRobin *roundRobin
	instances map[string]string
	pongs chan Pong
	mu sync.Mutex
}

func (dns *DNS) startPinging(ctx context.Context) {
	if dns.broker == nil {
		return
	}
	dnsLogger.Debug("dns pinging started")
	timeout := time.After(0)
	for {
		select {
		case <-timeout:
			fmt.Println("in timeout!!!")
			id := uuid.New().String()
			dns.broker.Emit(notify.BrokerEvent{
				Data:  id,
				Event: PingEvent,
			})
			deadline := time.After(dns.pongDeadline)

			replies := map[string]string{}
			fmt.Println("waiting for pongs")

			for {
				f := false
				select {
				case <-ctx.Done():
					return
				case pong := <-dns.pongs:
					fmt.Println("PONG IS HERE: ", pong)
					if pong.ID == id {
						replies[pong.BrokerID] = pong.IP
						dns.Add(pong.BrokerID, pong.IP)
					}
				case <-deadline:
					f = true
				}
				if f {
					break
				}
			}

			dnsLogger.Debug("replies")
			for rid := range replies {
				dnsLogger.Debug("[reply]", rid)
			}

			dns.mu.Lock()
			for brokerID := range dns.instances {

				if _, ok := replies[brokerID]; !ok {
					delete(dns.instances, brokerID)
					dns.roundRobin.del(brokerID)
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
			fmt.Println("PONG EVENT")
			if str, ok := e.Data.(string); ok {
				if err := json.Unmarshal([]byte(str), &pong); err == nil {
					pong.BrokerID = e.BrokerID
					select {
					case dns.pongs <- pong:
						fmt.Println("enqueued pong")
					}
				}
			}
		}
	})
	defer handle.Close()

	<-ctx.Done()
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

	dnsLogger.Debugf("failed to get address of %s", id)
	return "", errors.New("failed to get instance")
}

func (dns *DNS) Add(id, ip string) {
	dns.mu.Lock()
	defer dns.mu.Unlock()
	dnsLogger.Debugf("adding instance %s", id)
	dns.roundRobin.add(id)
	dns.instances[id] = ip
}

func (dns *DNS) Del(id string) {
	dns.mu.Lock()
	defer dns.mu.Unlock()
	dnsLogger.Debugf("deleting instance %s", id)
	delete(dns.instances, id)
	dns.roundRobin.del(id)
}

// Starts the dns server (non-blocking)
func (dns *DNS) Start(ctx context.Context) {
	go dns.startPinging(ctx)
	go dns.handleBrokerEvents(ctx)
}

func NewDNS(cfg Config) *DNS {
	cfg = cfg.validate()

	return &DNS{
		broker:       cfg.Broker,
		pongs:        make(chan Pong, cfg.PongsBufferSize),
		pingInterval: cfg.PingInterval,
		instances:    map[string]string{},
		roundRobin:   newRoundRobin(),
	}
}

func handleDNSPingEvent(b notify.Broker, ipProvider func() string) io.Closer {
	return b.Handle(func(e notify.BrokerEvent) {
		fmt.Println("received broker event:", e)
		if e.Event == PingEvent {
			fmt.Println("handling ping event...")
			ip := ipProvider()
			if id, ok := e.Data.(string); ok {
				pong := Pong{
					ID:       id,
					IP:       ip,
				}

				data, err := json.Marshal(pong)
				if err != nil {
					fmt.Println("failed to marshal pong:", err)
					return
				}
				fmt.Println("emitting pong event")
				b.Emit(notify.BrokerEvent{
					Data:  string(data),
					Event: PongEvent,
				})
			} else {
				fmt.Println("wefojwepofweopfkwef")
			}
		}
	})
}
