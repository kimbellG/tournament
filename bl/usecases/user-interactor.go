package usecases

import (
	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/models"
)

type UserInteractor struct {
	UserRepo UserRepository
}

func (ui *UserInteractor) Save(user *models.User) (uuid.UUID, error) {
	id, err := ui.UserRepo.Save(user)
	if err != nil {
		return id, kerror.Errorf(err, "repository")
	}

	return id, nil
}

func (ui *UserInteractor) GetByID(id uuid.UUID) (*models.User, error) {
	user, err := ui.UserRepo.GetById(id)
	if err != nil {
		return nil, kerror.Errorf(err, "repository")
	}

	return user, nil
}

func (ui *UserInteractor) DeleteByID(id uuid.UUID) error {
	if err := ui.UserRepo.DeleteById(id); err != nil {
		return kerror.Errorf(err, "repository")
	}

	return nil
}

func (ui *UserInteractor) UpdateBalance(id uuid.UUID, addend float64) error {
	if err := ui.UserRepo.SumToBalance(id, addend); err != nil {
		return kerror.Errorf(err, "repository")
	}

	return nil
}
