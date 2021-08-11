package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentInteractor struct {
	repo      TournamentRepository
	store     tx.Transactioner
	userCases UserUsecase
}

func (tu *TournamentInteractor) Create(tournament *models.Tournament) (uuid.UUID, error) {
	var id uuid.UUID

	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		id, err = tu.repo.Insert(store, tournament)
		if err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return id, kerror.Errorf(err, "transaction")
	}

	return id, nil

}

func (tu *TournamentInteractor) GetByID(id uuid.UUID) (*models.Tournament, error) {
	var tournament *models.Tournament

	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		tournament, err = tu.repo.SelectByID(store, id)
		if err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return nil, kerror.Errorf(err, "execution transaction")
	}

	return tournament, nil
}

func (tu *TournamentInteractor) Join(tournamentID uuid.UUID, userID uuid.UUID) error {
	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		isActiveTournament, err := tu.isActiveTournament(store, tournamentID)
		if err != nil {
			return kerror.Errorf(err, "check status of tournament")
		}

		if !isActiveTournament {
			return kerror.Newf(kerror.BadRequest, "tournament isn't active")
		}

		deposit, err := tu.getDeposit(store, tournamentID)
		if err != nil {
			return kerror.Errorf(err, "getting deposit")
		}

		// TODO: use repository function with transaction
		if err := tu.userCases.SumToBalance(userID, -deposit); err != nil {
			return kerror.Errorf(err, "subtraction from the balance")
		}

		if err := tu.repo.AddToPrize(store, tournamentID, deposit); err != nil {
			return kerror.Errorf(err, "adding to prize of tournament")
		}

		if err := tu.repo.InsertUserToTournament(store, tournamentID, userID); err != nil {
			return kerror.Errorf(err, "adding user to tournament")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}

func (tu *TournamentInteractor) getDeposit(store tx.DBTX, tournamentID uuid.UUID) (float64, error) {
	tournament, err := tu.repo.SelectByID(store, tournamentID)
	if err != nil {
		return -1, kerror.Errorf(err, "get tournament")
	}

	return tournament.Deposit, nil
}

func (tu *TournamentInteractor) isActiveTournament(store tx.DBTX, tournamentID uuid.UUID) (bool, error) {
	status, err := tu.getStatus(store, tournamentID)
	if err != nil {
		return false, kerror.Errorf(err, "get status of tournament")
	}

	return status == models.Active, nil
}

func (tu *TournamentInteractor) getStatus(store tx.DBTX, tournamentID uuid.UUID) (models.TournamentStatus, error) {
	tournament, err := tu.repo.SelectByID(store, tournamentID)
	if err != nil {
		return "", kerror.Errorf(err, "get tournament")
	}

	return tournament.Status, nil
}

func (tu *TournamentInteractor) Finish(id uuid.UUID) error {
	err := tu.store.WithTransaction(func(store tx.DBTX) error {

		isActiveTournament, err := tu.isActiveTournament(store, id)
		if err != nil {
			return kerror.Errorf(err, "check status of tournament")
		}

		if !isActiveTournament {
			return kerror.Newf(kerror.BadRequest, "tournament isn't active")
		}

		prize, err := tu.getPrize(store, id)
		if err != nil {
			return kerror.Errorf(err, "get prize")
		}

		winner, err := tu.generateWinner(store, id)
		if err != nil {
			return kerror.Errorf(err, "generate winner")
		}

		// TODO: use repository function with transaction
		if err := tu.userCases.SumToBalance(winner.ID, prize); err != nil {
			return kerror.Errorf(err, "add prize to winner's balance")
		}

		if err := tu.repo.UpdateStatus(store, id, models.Finish); err != nil {
			return kerror.Errorf(err, "change status")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}

func (tu *TournamentInteractor) generateWinner(store tx.DBTX, tournamentID uuid.UUID) (*models.User, error) {
	winner, err := tu.repo.SelectRandomUserOfTournament(store, tournamentID)
	if err != nil {
		return nil, kerror.Errorf(err, "get random user")
	}

	return winner, nil
}

func (tu *TournamentInteractor) getPrize(store tx.DBTX, tournamentID uuid.UUID) (float64, error) {
	tournament, err := tu.repo.SelectByID(store, tournamentID)
	if err != nil {
		return -1, kerror.Errorf(err, "get tournament")
	}

	return tournament.Prize, nil
}

func (tu *TournamentInteractor) Cancel(id uuid.UUID) error {
	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		isActiveTournament, err := tu.isActiveTournament(store, id)
		if err != nil {
			return kerror.Errorf(err, "check status of tournament")
		}

		if !isActiveTournament {
			return kerror.Newf(kerror.BadRequest, "tournament isn't active")
		}

		if err := tu.repo.RefundDepositToUsers(store, id); err != nil {
			return kerror.Errorf(err, "add deposit to user balance")
		}

		if err := tu.repo.UpdateStatus(store, id, models.Cancel); err != nil {
			return kerror.Errorf(err, "change status")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}
