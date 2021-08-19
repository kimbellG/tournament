package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament/core/debugutil"
	"github.com/kimbellG/tournament/core/models"
	"github.com/kimbellG/tournament/core/tx"
)

type UserRepository struct{}

func (u *UserRepository) Insert(ctx context.Context, store tx.DBTX, user *models.User) (uuid.UUID, error) {
	const query = `
		INSERT INTO Users(name, balance) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	inStmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return id, kerror.Newf(kerror.SQLPrepareStatementError, "prepare: %v", err)
	}
	defer debugutil.Close(inStmt)

	if err := inStmt.QueryRowContext(ctx, user.Name, user.Balance).Scan(&id); err != nil {
		return id, kerror.Newf(kerror.SQLConstraintError, "scan: %v", err)
	}

	return id, nil
}

func (u *UserRepository) SelectByID(ctx context.Context, store tx.DBTX, id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT * FROM Users WHERE id = $1;	
	`
	user := &models.User{}

	selectStmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(selectStmt)

	if err := selectStmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.Newf(kerror.UserDoesntExists, "no user with id %v: %v", id, err)
		}

		return nil, kerror.Newf(kerror.SQLScanError, "query: %v", err)
	}

	return user, nil
}

func (u *UserRepository) DeleteByID(ctx context.Context, store tx.DBTX, id uuid.UUID) error {
	const query = `
		DELETE FROM Users WHERE id = $1;	
	`

	deleteStmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare for deleting user: %v", err)
	}
	defer debugutil.Close(deleteStmt)

	if _, err := deleteStmt.ExecContext(ctx, id); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec deleting user: %v", err)
	}

	return nil
}

func (u *UserRepository) UpdateBalanceBySum(ctx context.Context, store tx.DBTX, id uuid.UUID, d float64) error {
	const query = `
		UPDATE Users
		SET balance = balance + $1
		WHERE id = $2
	`

	updateStmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare for update balance for %v(addend: %v): %v", id, d, err)
	}
	defer debugutil.Close(updateStmt)

	if _, err := updateStmt.ExecContext(ctx, d, id); err != nil {
		return kerror.Newf(kerror.SQLConstraintError, "updating balance for %v(addend: %v): %v", id, d, err)
	}

	return nil
}
