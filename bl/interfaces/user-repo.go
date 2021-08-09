package interfaces

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
)

type DB interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type UserRepository struct {
	db DB
}

func (u *UserRepository) Save(user *models.User) (uuid.UUID, error) {
	const query = `
		INSERT INTO Users(name, balance) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	inStmt, err := u.db.Prepare(query)
	if err != nil {
		return id, kerror.New(fmt.Errorf("prepare: %v", err), kerror.IntervalServerError)
	}

	if err := inStmt.QueryRow(user.Name, user.Balance).Scan(&id); err != nil {
		return id, kerror.New(fmt.Errorf("scan: %v", err), kerror.BadRequest)
	}

	return id, nil
}

func (u *UserRepository) GetById(id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT * FROM Users WHERE id = $1;	
	`
	user := &models.User{}

	selectStmt, err := u.db.Prepare(query)
	if err != nil {
		return nil, kerror.New(fmt.Errorf("prepare stmt: %v", err), kerror.IntervalServerError)
	}

	if err := selectStmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.New(kerror.Errorf(err, "no user with id $v", id), kerror.InvalidID)
		}

		return nil, kerror.New(kerror.Errorf(err, "query"), kerror.BadRequest)
	}

	return user, nil
}

func (u *UserRepository) DeleteById(id uuid.UUID) error {
	const query = `
		DELETE FROM Users WHERE id = $1;	
	`

	deleteStmt, err := u.db.Prepare(query)
	if err != nil {
		return kerror.New(fmt.Errorf("prepare for deleting user: %v", err), kerror.IntervalServerError)
	}

	if _, err := deleteStmt.Exec(id); err != nil {
		return kerror.New(fmt.Errorf("exec deleting user: %v", err), kerror.IntervalServerError)
	}

	return nil
}

func (u *UserRepository) SumToBalance(id uuid.UUID, addend float64) error {
	const query = `
		UPDATE Users
		SET balance = balance + $1
		WHERE id = $2
	`

	updateStmt, err := u.db.Prepare(query)
	if err != nil {
		return kerror.New(fmt.Errorf("prepare for update balance for %v(addend: %v): %v", id, addend, err), kerror.IntervalServerError)
	}

	if _, err := updateStmt.Exec(addend, id); err != nil {
		return kerror.New(fmt.Errorf("updating balance for %v(addend: %v): %v", id, addend, err), kerror.IntervalServerError)
	}

	return nil
}
