package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentRepository interface {
	Insert(repo tx.DBTX, tournament *models.Tournament) (uuid.UUID, error)
	InsertUserToTournament(repo tx.DBTX, tournamentID, userID uuid.UUID) error

	SelectByID(repo tx.DBTX, id uuid.UUID) (*models.Tournament, error)
	SelectRandomUserOfTournament(repo tx.DBTX, tournamentID uuid.UUID) (*models.User, error)

	AddToPrize(repo tx.DBTX, ID uuid.UUID, end float64) error
	RefundDepositToUsers(repo tx.DBTX, tournamentID uuid.UUID) error
	SetWinner(repo tx.DBTX, tournamentID, userID uuid.UUID) error

	UpdateStatus(repo tx.DBTX, tournamentID uuid.UUID, newStatus models.TournamentStatus) error
}
