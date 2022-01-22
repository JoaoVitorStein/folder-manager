package manager

import "github.com/jmoiron/sqlx"

type manager struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *manager {
	return &manager{
		db: db,
	}
}
