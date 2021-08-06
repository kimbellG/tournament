package interfaces

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
)

type UserRepository struct {
	db *sql.DB
}

func (u *UserRepository) Save(user *models.User) (uuid.UUID, error) {
	const query = `
		INSERT INTO Users(name, balance) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID
	retID := u.db.QueryRow(query, user.Name, user.Balance)
	if err := retID.Scan(id); err != nil {
		return id, kerror.New(fmt.Errorf("scan: %v", err), kerror.InvalidID)
	}

	return id, nil
}
