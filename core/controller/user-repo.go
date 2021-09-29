package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/tournament/core/models"
	"github.com/kimbellG/tournament/core/tx"
)

type UserRepository interface {
	Insert(ctx context.Context, store tx.DBTX, user *models.User) (uuid.UUID, error)
	SelectByID(ctx context.Context, store tx.DBTX, id uuid.UUID) (*models.User, error)
	SelectByName(ctx context.Context, store tx.DBTX, username string) (*models.User, error)
	DeleteByID(ctx context.Context, store tx.DBTX, id uuid.UUID) error
	UpdateBalanceBySum(ctx context.Context, store tx.DBTX, id uuid.UUID, d float64) error
}
