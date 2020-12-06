package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/cluster/broker"
	"github.com/RomanIschenko/notify/cluster/internal/protocol"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/sirupsen/logrus"
	"time"
)

var logger = logrus.WithField("source", "cluster")

type BrokerMetaInfo struct {}

func RegisterApp(ctx context.Context, b broker.Broker, app *notify.App) error {
	if app == nil {
		return errors.New("can not register a nil app")
	}
	if b == nil {
		return errors.New(fmt.Sprintf("can't register an app (%s) in a nil broker", app.ID()))
	}
	brokerManager := protocol.NewManager(app.ID(), b)
	b.Subscribe([]string{brokerManager.ChannelFormatter().AppChannel()}, time.Now().UnixNano())
	b.Handle(func(e broker.Event) {
		switch e.Name {
		case pubsub.PublishEvent:
			opts, err := brokerManager.ReadPublishOptions(e)
			if err != nil {
				logger.Error(err)
			}
			opts.MetaInfo = BrokerMetaInfo{}
			app.Publish(opts)
		case pubsub.SubscribeEvent:
			opts, err := brokerManager.ReadSubscribeOptions(e)
			if err != nil {
				logger.Error(err)
			}
			opts.MetaInfo = BrokerMetaInfo{}
			app.Subscribe(opts)
		case pubsub.UnsubscribeEvent:
			opts, err := brokerManager.ReadUnsubscribeOptions(e)
			if err != nil {
				logger.Error(err)
			}
			opts.MetaInfo = BrokerMetaInfo{}
			app.Unsubscribe(opts)
		}
	})
	appEvents := app.Events(ctx)
	appEvents.OnPublish(func(p pubsub.Publication, opts pubsub.PublishOptions) {
		if _, ok := opts.MetaInfo.(BrokerMetaInfo); ok {
			return
		}
		if err := brokerManager.WritePubsubOptions(opts); err != nil {
			logger.Error(err)
		}
	})
	appEvents.OnSubscribe(func(opts pubsub.SubscribeOptions, log changelog.Log) {
		if _, ok := opts.MetaInfo.(BrokerMetaInfo); ok {
			return
		}
		if err := brokerManager.WritePubsubOptions(opts); err != nil {
			logger.Error(err)
		}
	})
	appEvents.OnUnsubscribe(func(opts pubsub.UnsubscribeOptions, log changelog.Log) {
		if _, ok := opts.MetaInfo.(BrokerMetaInfo); ok {
			return
		}
		if err := brokerManager.WritePubsubOptions(opts); err != nil {
			logger.Error(err)
		}
	})
	appEvents.OnChange(func(log changelog.Log) {
		f := brokerManager.ChannelFormatter()
		var subs, unsubs []string

		for _, s := range log.TopicsUp {
			subs = append(subs, f.FmtTopic(s))
		}
		for _, s := range log.UsersUp {
			subs = append(subs, f.FmtUser(s))
		}
		for _, s := range log.ClientsUp {
			subs = append(subs, f.FmtClient(s))
		}

		for _, s := range log.TopicsDown {
			unsubs = append(unsubs, f.FmtTopic(s))
		}
		for _, s := range log.UsersDown {
			unsubs = append(unsubs, f.FmtUser(s))
		}
		for _, s := range log.ClientsDown {
			unsubs = append(unsubs, f.FmtClient(s))
		}
		b.Subscribe(subs, log.Time)
		b.Unsubscribe(unsubs, log.Time)
	})
	b.Subscribe(brokerManager.ChannelFormatter().FmtChannels(app.Clients(), app.Users(), app.Topics()), time.Now().UnixNano())
	return nil
}