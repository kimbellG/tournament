package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type UserInteractor struct {
	UserRepo UserRepository
	store    tx.Transactioner
}

func (ui *UserInteractor) Save(user *models.User) (uuid.UUID, error) {
	var id uuid.UUID

	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		id, err = ui.UserRepo.Insert(store, user)
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

func (ui *UserInteractor) GetByID(id uuid.UUID) (*models.User, error) {
	var user *models.User

	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		user, err = ui.UserRepo.SelectByID(store, id)
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

func (ui *UserInteractor) DeleteByID(id uuid.UUID) error {
	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		if err := ui.UserRepo.DeleteByID(store, id); err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}

func (ui *UserInteractor) UpdateBalance(id uuid.UUID, addend float64) error {
	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		if err := ui.UserRepo.UpdateBalanceBySum(store, id, addend); err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return kerror.Errorf(err, "execution transaction")
	}

	return nil
}
