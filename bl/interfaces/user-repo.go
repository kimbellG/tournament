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
		// TODO: Upgrade kerror and set status: IntervalServerError
		return id, kerror.New(fmt.Errorf("prepare: %v", err), kerror.BadRequest)
	}

	if err := inStmt.QueryRow(user.Name, user.Balance).Scan(id); err != nil {
		return id, kerror.New(fmt.Errorf("scan: %v", err), kerror.BadRequest)
	}

	return id, nil
}
