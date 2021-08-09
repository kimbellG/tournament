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
