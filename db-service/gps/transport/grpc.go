package transport

import (
	"context"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/mvromanchak/emil-development/db-service/gps"
	pr "github.com/mvromanchak/emil-development/db-service/protorepo"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type grpcServer struct {
	addGPS grpctransport.Handler
	logger log.Logger
}

// AddGRPCHandler adds a gRPC handler for groups endpoint.
func AddGRPCHandler(server *grpc.Server, endpoints gps.Endpoints, log log.Logger) {
	options := UnaryClientGRPCOptions()
	h := &grpcServer{
		addGPS: grpctransport.NewServer(
			endpoints[gps.CreateGPSEndpoint],
			decodeGPSRequest,
			encodeGPSResponse,
			options...,
		),
		logger: log,
	}
	pr.RegisterGpsAPIServer(server, h)
}

// ListGroups returns a filtered list of groups.
func (s *grpcServer) ListGps(ctx context.Context, req *pr.GpsRequest) (*pr.GpsResponse, error) {
	_, out, err := s.addGPS.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, errors.New("addGPS: no response returned")
	}
	return out.(*pr.GpsResponse), err
}

func decodeGPSRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	v, ok := grpcReq.(*pr.GpsRequest)
	if !ok {
		panic("type GPSRequest")
	}
	req := gps.GPSRequest{
		DeviceId: v.DeviceId,
	}
	return req, nil
}

func encodeGPSResponse(_ context.Context, resp interface{}) (interface{}, error) {
	_, ok := resp.(gps.GPSResponce)
	if !ok {
		panic("type assertion BlockRequest failed")
	}
	prr := pr.GpsResponse{Ok: true}
	return &prr, nil
}

// UnaryClientGRPCOptions is a set of the standard options applied to grpc client.
func UnaryClientGRPCOptions() []grpctransport.ServerOption {
	return []grpctransport.ServerOption{}
}
