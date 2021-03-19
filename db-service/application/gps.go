package application

import (
	"context"
	"github.com/mvromanchak/emil-development/db-service/entities"
	"github.com/mvromanchak/emil-development/db-service/gps"
)

type GPSRepository interface {
	AddGPS(ctx context.Context, gps *entities.GPSData) error
}

// NewGroupsService instantiate new groups service
func NewGroupsService(db GPSRepository) gps.Service {
	return &groupsService{
		db: db,
	}
}

// groupsService is the service of the application level for groups
type groupsService struct {
	db GPSRepository
}

func (s groupsService) AddGPS(ctx context.Context, gps *entities.GPSData) (err error) {
	return s.db.AddGPS(ctx, gps)
}
