package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
)

type UserController interface {
	Save(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	UpdateBalance(ctx context.Context, id uuid.UUID, addend float64) error
}
