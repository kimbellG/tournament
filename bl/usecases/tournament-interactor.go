package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
)

type TournamentInteractor struct {
	repo      TournamentRepository
	userCases UserUsecase
}

func (tu *TournamentInteractor) Create(tournament *models.Tournament) (uuid.UUID, error) {
	id, err := tu.repo.Create(tournament)
	if err != nil {
		return id, kerror.Errorf(err, "repo")
	}

	return id, nil
}

func (tu *TournamentInteractor) GetByID(id uuid.UUID) (*models.Tournament, error) {
	tournament, err := tu.repo.GetByID(id)
	if err != nil {
		return nil, kerror.Errorf(err, "repo")
	}

	return tournament, nil
}

func (tu *TournamentInteractor) Join(tournamentID uuid.UUID, userID uuid.UUID) error {
	isActiveTournament, err := tu.isActiveTournament(tournamentID)
	if err != nil {
		return kerror.Errorf(err, "check status of tournament")
	}

	if !isActiveTournament {
		return kerror.Newf(kerror.BadRequest, "tournament isn't active")
	}

	deposit, err := tu.getDeposit(tournamentID)
	if err != nil {
		return kerror.Errorf(err, "getting deposit")
	}

	if err := tu.userCases.SumToBalance(userID, -deposit); err != nil {
		return kerror.Errorf(err, "subtraction from the balance")
	}

	if err := tu.repo.AddToPrize(tournamentID, deposit); err != nil {
		return kerror.Errorf(err, "adding to prize of tournament")
	}

	if err := tu.repo.AddUserToTournament(tournamentID, userID); err != nil {
		return kerror.Errorf(err, "adding user to tournament")
	}

	return nil
}

func (tu *TournamentInteractor) getDeposit(tournamentID uuid.UUID) (float64, error) {
	tournament, err := tu.GetByID(tournamentID)
	if err != nil {
		return -1, kerror.Errorf(err, "get tournament")
	}

	return tournament.Deposit, nil
}

func (tu *TournamentInteractor) isActiveTournament(tournamentID uuid.UUID) (bool, error) {
	status, err := tu.getStatus(tournamentID)
	if err != nil {
		return false, kerror.Errorf(err, "get status of tournament")
	}

	return status == models.Active, nil
}

func (tu *TournamentInteractor) getStatus(tournamentID uuid.UUID) (models.TournamentStatus, error) {
	tournament, err := tu.GetByID(tournamentID)
	if err != nil {
		return "", kerror.Errorf(err, "get tournament")
	}

	return tournament.Status, nil
}

func (tu *TournamentInteractor) Finish(id uuid.UUID) error {
	isActiveTournament, err := tu.isActiveTournament(id)
	if err != nil {
		return kerror.Errorf(err, "check status of tournament")
	}

	if !isActiveTournament {
		return kerror.Newf(kerror.BadRequest, "tournament isn't active")
	}

	prize, err := tu.getPrize(id)
	if err != nil {
		return kerror.Errorf(err, "get prize")
	}

	winner, err := tu.generateWinner(id)
	if err != nil {
		return kerror.Errorf(err, "generate winner")
	}

	if err := tu.userCases.SumToBalance(winner.ID, prize); err != nil {
		return kerror.Errorf(err, "add prize to winner's balance")
	}

	if err := tu.repo.ChangeStatus(id, models.Finish); err != nil {
		return kerror.Errorf(err, "change status")
	}

	return nil
}

func (tu *TournamentInteractor) generateWinner(tournamentID uuid.UUID) (*models.User, error) {
	winner, err := tu.repo.GetRandomUserOfTournament(tournamentID)
	if err != nil {
		return nil, kerror.Errorf(err, "get random user")
	}

	return winner, nil
}

func (tu *TournamentInteractor) getPrize(tournamentID uuid.UUID) (float64, error) {
	tournament, err := tu.repo.GetByID(tournamentID)
	if err != nil {
		return -1, kerror.Errorf(err, "get tournament")
	}

	return tournament.Prize, nil
}

func (tu *TournamentInteractor) Cancel(id uuid.UUID) error {
	isActiveTournament, err := tu.isActiveTournament(id)
	if err != nil {
		return kerror.Errorf(err, "check status of tournament")
	}

	if !isActiveTournament {
		return kerror.Newf(kerror.BadRequest, "tournament isn't active")
	}

	if err := tu.repo.AddDepositToUsersOfTournament(id); err != nil {
		return kerror.Errorf(err, "add deposit to user balance")
	}

	if err := tu.repo.ChangeStatus(id, models.Cancel); err != nil {
		return kerror.Errorf(err, "change status")
	}

	return nil
}
