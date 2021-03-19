package repository

import "github.com/jmoiron/sqlx"

// Repo is the PostgreSQL storage for the product data
type Repo struct {
	dbconn *sqlx.DB
}

// NewRepository...
func NewRepository(dbconn *sqlx.DB) (r *Repo) {
	r = &Repo{
		dbconn: dbconn,
	}
	return
}

// Close closes prepared statements of the repository.
func (r *Repo) Close() error {
	return nil
}
