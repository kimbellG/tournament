package interfaces

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/debugutil"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type UserRepository struct{}

func (u *UserRepository) Save(store tx.DBTX, user *models.User) (uuid.UUID, error) {
	const query = `
		INSERT INTO Users(name, balance) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	inStmt, err := store.Prepare(query)
	if err != nil {
		return id, kerror.Newf(kerror.SQLPrepareStatementError, "prepare: %v", err)
	}
	defer debugutil.Close(inStmt)

	if err := inStmt.QueryRow(user.Name, user.Balance).Scan(&id); err != nil {
		return id, kerror.Newf(kerror.SQLConstraintError, "scan: %v", err)
	}

	return id, nil
}

func (u *UserRepository) GetByID(store tx.DBTX, id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT * FROM Users WHERE id = $1;	
	`
	user := &models.User{}

	selectStmt, err := store.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(selectStmt)

	if err := selectStmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.Newf(kerror.UserDoesntExists, "no user with id %v: %v", id, err)
		}

		return nil, kerror.Newf(kerror.SQLScanError, "query: %v", err)
	}

	return user, nil
}

func (u *UserRepository) DeleteByID(store tx.DBTX, id uuid.UUID) error {
	const query = `
		DELETE FROM Users WHERE id = $1;	
	`

	deleteStmt, err := store.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare for deleting user: %v", err)
	}
	defer debugutil.Close(deleteStmt)

	if _, err := deleteStmt.Exec(id); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec deleting user: %v", err)
	}

	return nil
}

func (u *UserRepository) SumToBalance(store tx.DBTX, id uuid.UUID, addend float64) error {
	const query = `
		UPDATE Users
		SET balance = balance + $1
		WHERE id = $2
	`

	updateStmt, err := store.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare for update balance for %v(addend: %v): %v", id, addend, err)
	}
	defer debugutil.Close(updateStmt)

	if _, err := updateStmt.Exec(addend, id); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "updating balance for %v(addend: %v): %v", id, addend, err)
	}

	return nil
}
