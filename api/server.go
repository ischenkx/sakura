package api

import (
	"context"
	"encoding/json"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/api/dns"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/google/uuid"
	"net/http"
)

type Config struct {
	APIKey string
	DNS  dns.Config
	Auth notify.Auth
	IPS []string
}

type Input struct {
	APIKey string
	ClientID pubsub.ClientID
}

type Output struct {
	IP string
	AuthData string
}

type Server struct {
	dns *dns.DNS
	auth notify.Auth
	apiKey string
}

func (s *Server) Start(ctx context.Context) {
	s.dns.Start(ctx)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(403)
		return
	}
	if s.apiKey != input.APIKey && s.apiKey != "" {
		w.WriteHeader(403)
		return
	}
	ip, err := s.dns.Next()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	authData := string(input.ClientID)
	if s.auth != nil {
		authData, err = s.auth.Register(input.ClientID)
		if err != nil {
			w.WriteHeader(403)
			return
		}
	}
	json.NewEncoder(w).Encode(Output{
		IP:       ip,
		AuthData: authData,
	})
}

func NewServer(config Config) *Server {
	d := dns.New(config.DNS)
	for _, ip := range config.IPS {
		d.Add(uuid.New().String(), ip)
	}

	return &Server{
		dns: 	d,
		auth:   config.Auth,
		apiKey: config.APIKey,
	}
}