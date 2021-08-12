package controller

import (
	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type UserRepository interface {
	Insert(store tx.DBTX, user *models.User) (uuid.UUID, error)
	SelectByID(store tx.DBTX, id uuid.UUID) (*models.User, error)
	DeleteByID(store tx.DBTX, id uuid.UUID) error
	UpdateBalanceBySum(store tx.DBTX, id uuid.UUID, d float64) error
}
