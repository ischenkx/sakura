package api

import (
	"errors"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/cluster/broker"
	"github.com/RomanIschenko/notify/cluster/internal/protocol"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/client_id"
)

type Config struct {
	AppID string
	Auth notify.Auth
	Broker broker.Broker
}

type Api struct {
	appId     string
	auth      notify.Auth
	brokerManager *protocol.Manager
}

func (api *Api) Publish(opts pubsub.PublishOptions) error {
	if api.brokerManager == nil {
		return errors.New("no broker provided")
	}
	return api.brokerManager.WritePubsubOptions(opts)
}

func (api *Api) Subscribe(opts pubsub.SubscribeOptions) error {
	if api.brokerManager == nil {
		return errors.New("no broker provided")
	}
	return api.brokerManager.WritePubsubOptions(opts)
}

func (api *Api) Unsubscribe(opts pubsub.UnsubscribeOptions) error {
	if api.brokerManager == nil {
		return errors.New("no broker provided")
	}
	return api.brokerManager.WritePubsubOptions(opts)
}

func (api *Api) Authorize(token string) (clientid.ID, error){
	return api.auth.Authorize(token)
}

func New(cfg Config) *Api {
	return &Api{
		appId:         cfg.AppID,
		auth:          cfg.Auth,
		brokerManager: protocol.NewManager(cfg.AppID, cfg.Broker),
	}
}