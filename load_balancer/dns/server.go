package dnslb

import (
	"context"
	sd "github.com/RomanIschenko/notify/load_balancer/dns/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var serverLogger = logrus.WithField("source", "dns_server")

type Server struct {
	lb *loadBalancer
	sd.UnimplementedDNSLoadBalancerServer
}

func (s *Server) GetAddress(context.Context, *sd.DNSRequest) (*sd.DNSResponse, error) {
	addr, err := s.lb.next()
	return &sd.DNSResponse{
		Address: addr,
	}, err
}

func (s *Server) Start(ctx context.Context, lis net.Listener, opts... grpc.ServerOption) {
	grpcServer := grpc.NewServer(opts...)

	sd.RegisterDNSLoadBalancerServer(grpcServer, s)

	s.lb.start(ctx)

	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			serverLogger.Error("grpc server failed to start:", err)
		}
	}()

	select {
	case <-ctx.Done():
	}
	grpcServer.GracefulStop()
}

func NewServer(config Config) *Server {
	return &Server{
		lb: newLoadBalancer(config),
	}
}