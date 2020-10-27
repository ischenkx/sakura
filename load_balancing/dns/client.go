package dnslb

import (
	"context"
	sd "github.com/RomanIschenko/notify/load_balancing/dns/pb"
	"google.golang.org/grpc"
)

type Client struct {
	client sd.DNSLoadBalancerClient
}

func (c Client) GetAddress(ctx context.Context, opts ...grpc.CallOption) (string, error) {
	res, err := c.client.GetAddress(ctx, &sd.DNSRequest{}, opts...)

	if err != nil {
		return "", err
	}

	return res.Address, nil
}

func NewClient(connInterface grpc.ClientConnInterface) Client {
	return Client{sd.NewDNSLoadBalancerClient(connInterface)}
}