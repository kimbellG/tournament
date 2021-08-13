package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type UserInteractor struct {
	UserRepo UserRepository
	store    tx.Store
}

func NewUserController(repo UserRepository, store tx.Store) UserController {
	return &UserInteractor{
		UserRepo: repo,
		store:    store,
	}
}

func (ui *UserInteractor) Save(ctx context.Context, user *models.User) (uuid.UUID, error) {
	var id uuid.UUID

	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		id, err = ui.UserRepo.Insert(ctx, store, user)
		if err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return id, kerror.Errorf(err, "execution transactive")
	}

	return id, nil
}

func (ui *UserInteractor) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user *models.User

	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		user, err = ui.UserRepo.SelectByID(ctx, store, id)
		if err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return nil, kerror.Errorf(err, "execution transactive")
	}

	return user, nil
}

func (ui *UserInteractor) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		if err := ui.UserRepo.DeleteByID(ctx, store, id); err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}

func (ui *UserInteractor) UpdateBalance(ctx context.Context, id uuid.UUID, addend float64) error {
	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		if err := ui.UserRepo.UpdateBalanceBySum(ctx, store, id, addend); err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}
