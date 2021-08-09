package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
)

type UserRepository interface {
	Save(user *models.User) (uuid.UUID, error)
	GetById(id uuid.UUID) (*models.User, error)
	DeleteById(id uuid.UUID) error
	SumToBalance(id uuid.UUID, addend float64) error
}
