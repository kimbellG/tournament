package controller

import (
	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
)

type UserUsecase interface {
	Save(user *models.User) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*models.User, error)
	DeleteByID(id uuid.UUID) error
	SumToBalance(id uuid.UUID, addend float64) error
}
