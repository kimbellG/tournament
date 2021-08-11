package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentRepository interface {
	Create(repo tx.DBTX, tournament *models.Tournament) (uuid.UUID, error)

	GetByID(repo tx.DBTX, id uuid.UUID) (*models.Tournament, error)
	GetRandomUserOfTournament(repo tx.DBTX, tournamentID uuid.UUID) (*models.User, error)

	AddUserToTournament(repo tx.DBTX, tournamentID, userID uuid.UUID) error
	AddToPrize(repo tx.DBTX, ID uuid.UUID, addend float64) error
	AddDepositToUsersOfTournament(repo tx.DBTX, tournamentID uuid.UUID) error
	AddWinner(repo tx.DBTX, tournamentID, userID uuid.UUID) error

	ChangeStatus(repo tx.DBTX, tournamentID uuid.UUID, newStatus models.TournamentStatus) error
}
