package db

import (
	"database/sql"
)

type Store interface {
	Querier
}

type Storage struct {
	db *sql.DB
	*Queries
}

func NewStorage(db *sql.DB) Store {
	return &Storage{
		db:      db,
		Queries: New(db),
	}
}
