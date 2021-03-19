package gps

import (
	"context"
	"github.com/mvromanchak/emil-development/db-service/entities"
)

// Service is the domain interface for the products service
type Service interface {
	AddGPS(ctx context.Context, gps *entities.GPSData) error
}

// NewGroupsService instantiate core domain of group service
func NewGroupsService(s Service) Service {
	return &service{
		groupService: s,
	}
}

type service struct {
	groupService Service
}

func (s *service) AddGPS(ctx context.Context, gps *entities.GPSData) error {
	return s.groupService.AddGPS(ctx, gps)
}
