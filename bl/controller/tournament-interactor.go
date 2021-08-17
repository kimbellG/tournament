package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentInteractor struct {
	repo     TournamentRepository
	store    tx.Store
	userRepo UserRepository
}

func NewTournamentController(repo TournamentRepository, userRepo UserRepository, store tx.Store) TournamentController {
	return &TournamentInteractor{
		repo:     repo,
		userRepo: userRepo,
		store:    store,
	}
}

func (tu *TournamentInteractor) Create(ctx context.Context, tournament *models.Tournament) (uuid.UUID, error) {
	var id uuid.UUID

	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		id, err = tu.repo.Insert(ctx, store, tournament)
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

func (tu *TournamentInteractor) GetByID(ctx context.Context, id uuid.UUID) (*models.Tournament, error) {
	var tournament *models.Tournament

	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		tournament, err = tu.repo.SelectByID(ctx, store, id)
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

func (tu *TournamentInteractor) Join(ctx context.Context, tournamentID uuid.UUID, userID uuid.UUID) error {
	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		isActiveTournament, err := tu.isActiveTournament(ctx, store, tournamentID)
		if err != nil {
			return kerror.Errorf(err, "check status of tournament")
		}

		if !isActiveTournament {
			return kerror.Newf(kerror.BadRequest, "tournament isn't active")
		}

		deposit, err := tu.getDeposit(ctx, store, tournamentID)
		if err != nil {
			return kerror.Errorf(err, "getting deposit")
		}

		if err := tu.userRepo.UpdateBalanceBySum(ctx, store, userID, -deposit); err != nil {
			return kerror.Errorf(err, "subtraction from the balance")
		}

		if err := tu.repo.AddToPrize(ctx, store, tournamentID, deposit); err != nil {
			return kerror.Errorf(err, "adding to prize of tournament")
		}

		if err := tu.repo.InsertUserToTournament(ctx, store, tournamentID, userID); err != nil {
			return kerror.Errorf(err, "adding user to tournament")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}

func (tu *TournamentInteractor) getDeposit(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) (float64, error) {
	tournament, err := tu.repo.SelectByID(ctx, store, tournamentID)
	if err != nil {
		return -1, kerror.Errorf(err, "get tournament")
	}

	return tournament.Deposit, nil
}

func (tu *TournamentInteractor) isActiveTournament(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) (bool, error) {
	status, err := tu.getStatus(ctx, store, tournamentID)
	if err != nil {
		return false, kerror.Errorf(err, "get status of tournament")
	}

	return status == models.Active, nil
}

func (tu *TournamentInteractor) getStatus(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) (models.TournamentStatus, error) {
	tournament, err := tu.repo.SelectByID(ctx, store, tournamentID)
	if err != nil {
		return "", kerror.Errorf(err, "get tournament")
	}

	return tournament.Status, nil
}

func (tu *TournamentInteractor) Finish(ctx context.Context, id uuid.UUID) error {
	err := tu.store.WithTransaction(func(store tx.DBTX) error {

		isActiveTournament, err := tu.isActiveTournament(ctx, store, id)
		if err != nil {
			return kerror.Errorf(err, "check status of tournament")
		}

		if !isActiveTournament {
			return kerror.Newf(kerror.BadRequest, "tournament isn't active")
		}

		prize, err := tu.getPrize(ctx, store, id)
		if err != nil {
			return kerror.Errorf(err, "get prize")
		}

		winner, err := tu.generateWinner(ctx, store, id)
		if err != nil {
			return kerror.Errorf(err, "generate winner")
		}

		if err := tu.userRepo.UpdateBalanceBySum(ctx, store, winner.ID, prize); err != nil {
			return kerror.Errorf(err, "add prize to winner's balance")
		}

		if err := tu.repo.SetWinner(ctx, store, id, winner.ID); err != nil {
			return kerror.Errorf(err, "set winner")
		}

		if err := tu.repo.UpdateStatus(ctx, store, id, models.Finish); err != nil {
			return kerror.Errorf(err, "change status")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}

func (tu *TournamentInteractor) generateWinner(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) (*models.User, error) {
	winner, err := tu.repo.SelectRandomUserOfTournament(ctx, store, tournamentID)
	if err != nil {
		return nil, kerror.Errorf(err, "get random user")
	}

	return winner, nil
}

func (tu *TournamentInteractor) getPrize(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) (float64, error) {
	tournament, err := tu.repo.SelectByID(ctx, store, tournamentID)
	if err != nil {
		return -1, kerror.Errorf(err, "get tournament")
	}

	return tournament.Prize, nil
}

func (tu *TournamentInteractor) Cancel(ctx context.Context, id uuid.UUID) error {
	err := tu.store.WithTransaction(func(store tx.DBTX) error {
		isActiveTournament, err := tu.isActiveTournament(ctx, store, id)
		if err != nil {
			return kerror.Errorf(err, "check status of tournament")
		}

		if !isActiveTournament {
			return kerror.Newf(kerror.BadRequest, "tournament isn't active")
		}

		if err := tu.repo.RefundDepositToUsers(ctx, store, id); err != nil {
			return kerror.Errorf(err, "add deposit to user balance")
		}

		if err := tu.repo.UpdateStatus(ctx, store, id, models.Cancel); err != nil {
			return kerror.Errorf(err, "change status")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}
