package usecases

import "github.com/google/uuid"

type TournamentRepository interface {
	Create(tournament *models.Tournament) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*models.Tournament, error)
	AddUserToTournament(tournamentID uuid.UUID, userID uuid.UUID) error
	ChangePrice(ID uuid.UUID, addend float64) error
	GetWinner(tournamentID uuid.UUID) (*models.Users, error)
	AddDepositToUsersOfTournament(tournamentID uuid.UUID) error
	// TODO: added this func after merge interfaces
	//	ChangeStatus(tournamentID uuid.UUID, newStatus )

}
