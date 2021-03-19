package application

import (
	"context"

	"github.com/mvromanchak/emil-development/api-service/entities"
	"github.com/mvromanchak/emil-development/api-service/gps/infrastructure"
	pr "github.com/mvromanchak/emil-development/api-service/protorepo"
	"google.golang.org/grpc"
)

type Client interface {
	SetGPS(ctx context.Context, data *entities.GPSData) (bool, error)
}

func NewGPSService(client infrastructure.GPSClient) gpsService {
	return gpsService{
		client: client,
	}
}

type gpsService struct {
	client infrastructure.GPSClient
}

func (i gpsService) SetGPS(ctx context.Context, data *entities.GPSData) (bool, error) {
	opts := []grpc.CallOption{}
	r := pr.GpsRequest{DeviceId: "test"}
	resp, err := i.client.ListGps(ctx, &r, opts...)
	if err != nil {
		return false, err
	}
	return resp.Ok, err
}
