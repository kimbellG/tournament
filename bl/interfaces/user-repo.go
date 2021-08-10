package interfaces

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/debugutil"
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

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Save(user *models.User) (uuid.UUID, error) {
	const query = `
		INSERT INTO Users(name, balance) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	inStmt, err := u.db.Prepare(query)
	if err != nil {
		return id, kerror.Newf(kerror.IntervalServerError, "prepare: %v", err)
	}
	defer debugutil.Close(inStmt)

	if err := inStmt.QueryRow(user.Name, user.Balance).Scan(&id); err != nil {
		return id, kerror.Newf(kerror.BadRequest, "scan: %v", err)
	}

	return id, nil
}

func (u *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT * FROM Users WHERE id = $1;	
	`
	user := &models.User{}

	selectStmt, err := u.db.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.IntervalServerError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(selectStmt)

	if err := selectStmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.Newf(kerror.InvalidID, "no user with id %v: %v", id, err)
		}

		return nil, kerror.Newf(kerror.BadRequest, "query: %v", err)
	}

	return user, nil
}

func (u *UserRepository) DeleteByID(id uuid.UUID) error {
	const query = `
		DELETE FROM Users WHERE id = $1;	
	`

	deleteStmt, err := u.db.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.IntervalServerError, "prepare for deleting user: %v", err)
	}
	defer debugutil.Close(deleteStmt)

	if _, err := deleteStmt.Exec(id); err != nil {
		return kerror.Newf(kerror.IntervalServerError, "exec deleting user: %v", err)
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
		return kerror.Newf(kerror.IntervalServerError, "prepare for update balance for %v(addend: %v): %v", id, addend, err)
	}
	defer debugutil.Close(updateStmt)

	if _, err := updateStmt.Exec(addend, id); err != nil {
		return kerror.Newf(kerror.IntervalServerError, "updating balance for %v(addend: %v): %v", id, addend, err)
	}

	return nil
}
