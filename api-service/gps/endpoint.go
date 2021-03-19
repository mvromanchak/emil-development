package gps

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/mvromanchak/emil-development/api-service/entities"
)

const (
	// CreateGPSEndpoint.
	CreateGPSEndpoint = "CreateGPSEndpoint"
)

// Endpoints contains the endpoints for the gps service functionality.
type Endpoints map[string]endpoint.Endpoint

// NewEndpoints creates new endpoints for gps service.
func NewEndpoints(svc Service) Endpoints {
	return Endpoints{
		CreateGPSEndpoint: makeCreateGPSEndpoint(svc),
	}
}

func makeCreateGPSEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		gpx, ok := request.(entities.GPSDataRequest)
		if !ok {
			panic("type assertion GPSDataRequestfailed")
		}
		// auth check.
		claims, err := gpx.Claims(ctx, "test", jwt.SigningMethodHS256)
		if err != nil {
			return nil, err
		}
		gpx.Authorize(gpx.DeviceID, claims)
		if err != nil {
			return nil, err
		}

		data := entities.GPSData{DeviceID: gpx.DeviceID}
		ok, err = s.SetGPS(ctx, &data)
		if err != nil {
			return nil, err
		}
		resp := entities.GPXResponse{OK: ok}
		return resp, nil
	}
}
