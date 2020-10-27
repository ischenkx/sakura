package api

import (
	"context"
	notifier "github.com/RomanIschenko/notify/api/pb"
	"github.com/RomanIschenko/notify/pubsub"
	"google.golang.org/grpc"
)

type Client struct {
	appID string
	pbClient notifier.NotifyClient
}

func (c Client) Publish(ctx context.Context, options pubsub.PublishOptions, opts ...grpc.CallOption) error {
	grpcOpts := notifier.PublishRequest{
		Clients: options.Clients,
		Users:   options.Users,
		Topics:  options.Topics,
		Message: options.Payload.Data,
		Time:    options.Time,
		AppID:   c.appID,
	}
	_, err := c.pbClient.Publish(ctx, &grpcOpts, opts...)
	return err
}

func (c Client) Subscribe(ctx context.Context, options pubsub.SubscribeOptions, opts ...grpc.CallOption) error {
	grpcOpts := notifier.SubscribeRequest{
		Clients: options.Clients,
		Users:   options.Users,
		Topics:  options.Topics,
		Time:    options.Time,
		AppID:   c.appID,
	}
	_, err := c.pbClient.Subscribe(ctx, &grpcOpts, opts...)
	return err
}

func (c Client) Unsubscribe(ctx context.Context, options pubsub.UnsubscribeOptions, opts ...grpc.CallOption) error {
	grpcOpts := notifier.UnsubscribeRequest{
		Clients: options.Clients,
		Users:   options.Users,
		Topics:  options.Topics,
		All: 	 options.All,
		Time:    options.Time,
		AppID:   c.appID,
	}
	_, err := c.pbClient.Unsubscribe(ctx, &grpcOpts, opts...)
	return err
}

func (c Client) Authorize(ctx context.Context, id pubsub.ClientID, opts ...grpc.CallOption) (*notifier.AuthOutput, error){
	grpcOpts := notifier.AuthInput{
		ClientID: id.String(),
	}
	return c.pbClient.Authorize(ctx, &grpcOpts, opts...)
}

func NewClient(appID string, cc grpc.ClientConnInterface) Client {
	return Client{appID, notifier.NewNotifyClient(cc)}
}
