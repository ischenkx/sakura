package protocol

import (
	"encoding/json"
	"errors"
	"github.com/RomanIschenko/notify/cluster/broker"
	"github.com/RomanIschenko/notify/pubsub"
)

type Manager struct {
	broker broker.Broker
	fmt *ChannelFormatter
}

func (m *Manager) ReadPublishOptions(e broker.Event) (opts pubsub.PublishOptions, err error) {
	opts.Clients, opts.Users, opts.Topics = m.fmt.ParseChannel(e.Source)
	var obj pubOpts
	if err = json.Unmarshal(e.Data, &obj); err != nil {
		return
	}
	opts.Payload = obj.Data
	opts.Time = obj.Time
	return
}

func (m *Manager) ReadSubscribeOptions(e broker.Event) (opts pubsub.SubscribeOptions, err error) {
	opts.Clients, opts.Users, opts.Topics = m.fmt.ParseChannel(e.Source)
	var obj subOpts
	if err = json.Unmarshal(e.Data, &obj); err != nil {
		return
	}
	opts.Topics = obj.Topics
	opts.Time = obj.Time
	return
}

func (m *Manager) ReadUnsubscribeOptions(e broker.Event) (opts pubsub.UnsubscribeOptions, err error) {
	opts.Clients, opts.Users, opts.Topics = m.fmt.ParseChannel(e.Source)
	var obj unsubOpts
	if err = json.Unmarshal(e.Data, &obj); err != nil {
		return
	}
	opts.Topics = obj.Topics
	opts.Time = obj.Time
	opts.All = obj.All
	return
}


func (m *Manager) WritePubsubOptions(data interface{}) error {
	var event string
	var channels []string
	var obj interface{}

	switch opts := data.(type) {
	case pubsub.PublishOptions:
		channels = m.fmt.FmtChannels(opts.Clients, opts.Users, opts.Topics)
		obj = pubOpts{
			Data: opts.Payload,
			Time: opts.Time,
		}
		event = pubsub.PublishEvent
	case pubsub.SubscribeOptions:
		channels = m.fmt.FmtChannels(opts.Clients, opts.Users, nil)
		obj = subOpts{
			Topics: opts.Topics,
			Time:   opts.Time,
		}
		event = pubsub.SubscribeEvent
	case pubsub.UnsubscribeOptions:
		channels = m.fmt.FmtChannels(opts.Clients, opts.Users, nil)
		obj = unsubOpts{
			Topics: opts.Topics,
			All:    opts.All,
			Time:   opts.Time,
		}
		event = pubsub.UnsubscribeEvent
	default:
		return errors.New("no options provided")
	}

	bts, err := json.Marshal(obj)

	if err != nil {
		return err
	}

	return m.broker.Publish(channels, broker.NewEvent(event, bts))
}

func (m *Manager) ChannelFormatter() *ChannelFormatter {
	return m.fmt
}

func NewManager(appId string, broker broker.Broker) *Manager {
	return &Manager{
		broker: broker,
		fmt:    New(appId),
	}
}