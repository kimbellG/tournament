package tx

import (
	"context"
	"database/sql"

	"github.com/kimbellG/kerror"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TransactionFunction func(store DBTX) error

type Store interface {
	WithTransaction(fn TransactionFunction) error
}

type TxStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &TxStore{
		db: db,
	}
}

func (s *TxStore) WithTransaction(fn TransactionFunction) error {
	tx, err := s.db.Begin()
	if err != nil {
		return kerror.Newf(kerror.SQLTransactionBeginError, "begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return kerror.Errorf(err, "transaction roolback error: %w; transaction error", rbErr)
		}

		return kerror.Errorf(err, "transaction function error")
	}

	if err := tx.Commit(); err != nil {
		return kerror.Newf(kerror.SQLTransactionCommitError, "transaction commit: %w", err)
	}

	return nil
}
