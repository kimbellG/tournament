package controller

import (
	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
)

type TournamentUsecase interface {
	Create(tournament *models.Tournament) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*models.Tournament, error)
	Join(tournamnetID uuid.UUID, userID uuid.UUID) error
	Finish(id uuid.UUID) error
	Cancel(id uuid.UUID) error
}
