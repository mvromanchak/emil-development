package infrastructure

import (
	"context"
	pr "github.com/mvromanchak/emil-development/api-service/protorepo"
	"google.golang.org/grpc"
)

type (
	// GPSClient is the stream client interface
	GPSClient interface {
		ListGps(ctx context.Context, in *pr.GpsRequest, opts ...grpc.CallOption) (*pr.GpsResponse, error)
	}
)

type prdClient struct {
	conn   *grpc.ClientConn
	client pr.GpsAPIClient
}

func (p prdClient) ListGps(ctx context.Context, in *pr.GpsRequest, opts ...grpc.CallOption) (*pr.GpsResponse, error) {
	r, err := p.client.ListGps(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func NewGPSGRPCClient(conn *grpc.ClientConn, opt ...func(ctx context.Context) context.Context) (GPSClient, error) {
	return &prdClient{
		conn:   conn,
		client: pr.NewGpsAPIClient(conn),
	}, nil
}
