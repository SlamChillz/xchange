package db

import (
	"fmt"
	"context"
	"database/sql"
)

type Store interface {
	Querier
	CreateCustomerTransaction(
		context.Context,
		CreateCustomerTransactionParams,
	) (Customer, error)
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

func (storage *Storage) executeTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", err, rbErr)
		}
		return fmt.Errorf("transaction err: %v", err)
	}
	return tx.Commit()
}
