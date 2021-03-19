package gps

import (
	"context"
	"github.com/mvromanchak/emil-development/api-service/entities"
)

type (
	// Service is the domain interface.
	Service interface {
		SetGPS(ctx context.Context, data *entities.GPSData) (bool, error)
	}
)

func NewGPSService(s Service) Service {
	return service{
		gpsService: s,
	}
}

type service struct {
	gpsService Service
}

func (s service) SetGPS(ctx context.Context, data *entities.GPSData) (bool, error) {
	return s.gpsService.SetGPS(ctx, data)
}
