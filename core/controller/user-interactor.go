package controller

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament/core/models"
	"github.com/kimbellG/tournament/core/tx"
	"golang.org/x/crypto/bcrypt"
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

func (ui *UserInteractor) Save(ctx context.Context, user *models.User) (*models.User, error) {
	created := &models.User{
		Name:     user.Name,
		Password: generatePassword(),
		Balance:  user.Balance,
	}

	hash, err := hashPassword(fmt.Sprintf("%x", sha256.Sum256([]byte(created.Password))))
	if err != nil {
		return nil, kerror.Newf(kerror.InternalServerError, "hashing password: %v", err)
	}
	user.Password = hash

	err = ui.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		created.ID, err = ui.UserRepo.Insert(ctx, store, user)
		if err != nil {
			return kerror.Errorf(err, "repository")
		}

		return nil
	})
	if err != nil {
		return nil, kerror.Errorf(err, "execution transactive")
	}

	return created, nil
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

func (ui *UserInteractor) Authorization(ctx context.Context, username, password string) (*models.User, error) {
	var user *models.User

	err := ui.store.WithTransaction(func(store tx.DBTX) error {
		var err error

		user, err = ui.UserRepo.SelectByName(ctx, store, username)
		if err != nil {
			return kerror.Errorf(err, "get user from database")
		}

		return nil
	})

	if err != nil {
		return nil, kerror.Errorf(err, "transactive database")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, kerror.Newf(kerror.IncorrectPassword, "compare password: %v", err)
	}

	return user, nil
}
