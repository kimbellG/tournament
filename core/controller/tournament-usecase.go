package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/tournament/core/models"
)

type TournamentController interface {
	Create(ctx context.Context, tournament *models.Tournament) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tournament, error)
	Join(ctx context.Context, tournamnetID uuid.UUID, userID uuid.UUID) error
	Finish(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
}
