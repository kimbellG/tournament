package tx

import (
	"database/sql"

	"github.com/kimbellG/kerror"
)

type DBTX interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type TransactionFunction func(store DBTX) error

type Transactioner interface {
	WithTransaction(fn TransactionFunction) error
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Transactioner {
	return &Store{
		db: db,
	}
}

func (s *Store) WithTransaction(fn TransactionFunction) error {
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
