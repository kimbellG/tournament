package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentRepository interface {
	Insert(ctx context.Context, repo tx.DBTX, tournament *models.Tournament) (uuid.UUID, error)
	InsertUserToTournament(ctx context.Context, repo tx.DBTX, tournamentID, userID uuid.UUID) error

	SelectByID(ctx context.Context, repo tx.DBTX, id uuid.UUID) (*models.Tournament, error)
	SelectRandomUserOfTournament(ctx context.Context, repo tx.DBTX, tournamentID uuid.UUID) (*models.User, error)

	AddToPrize(ctx context.Context, repo tx.DBTX, ID uuid.UUID, end float64) error
	RefundDepositToUsers(ctx context.Context, repo tx.DBTX, tournamentID uuid.UUID) error
	SetWinner(ctx context.Context, repo tx.DBTX, tournamentID, userID uuid.UUID) error

	UpdateStatus(ctx context.Context, repo tx.DBTX, tournamentID uuid.UUID, newStatus models.TournamentStatus) error
}
