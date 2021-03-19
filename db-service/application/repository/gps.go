package repository

import (
	"context"

	"github.com/mvromanchak/emil-development/db-service/entities"
)

func (r *Repo) AddGPS(ctx context.Context, gps *entities.GPSData) error {
	_, err := r.dbconn.NamedExec(`INSERT INTO signals (device_id) VALUES (:device_id)`,
		map[string]interface{}{
			"device_id": gps.DeviceId,
		})
	return err
}
