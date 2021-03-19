package gps

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mvromanchak/emil-development/db-service/entities"
)

const (
	// CreateGPSEndpoint.
	CreateGPSEndpoint = "CreateGPSEndpoint"
)

// Endpoints contains the endpoints for the service.
type Endpoints map[string]endpoint.Endpoint

// NewEndpoints creates new endpoints for the service.
func NewEndpoints(svc Service) Endpoints {
	return Endpoints{
		CreateGPSEndpoint: makeCreateGroupEndpoint(svc),
	}
}

type GPSRequest struct {
	DeviceId string
	GPS      []struct {
		Lat  string
		Lon  string
		Ele  string
		Time string
	}
}
type GPSResponce struct {
}

func makeCreateGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		v, ok := request.(GPSRequest)
		if !ok {
			panic("type assertion BlockRequest failed")
		}
		gd := entities.GPSData{
			DeviceId: v.DeviceId,
		}
		err = s.AddGPS(ctx, &gd)
		resp := GPSResponce{}
		return resp, err
	}
}
