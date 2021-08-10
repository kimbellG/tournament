package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	ttgrpc "github.com/kimbellG/tournament-bl/interfaces/grpc"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/usecases"
)

type ServiceController struct {
	userUsecase usecases.UserUsecase
}

func (sc *ServiceController) SaveUser(ctx context.Context, user *ttgrpc.User) (*ttgrpc.SaveResponse, error) {
	mUser, err := userGrpcToModels(user)
	if err != nil {
		return nil, kerror.Errorf(err, "marshaling user struct to models")
	}

	id, err := sc.userUsecase.Save(mUser)
	if err != nil {
		return nil, kerror.Errorf(err, "save user")
	}

	return &ttgrpc.SaveResponse{
		Id: id.String(),
	}, nil
}

func userGrpcToModels(gUser *ttgrpc.User) (*models.User, error) {
	mUser := &models.User{
		Name:    gUser.GetName(),
		Balance: gUser.GetBalance(),
	}

	id, err := uuid.Parse(gUser.GetID())
	if err != nil {
		return nil, kerror.Newf(kerror.InvalidID, "parser id: %v", err)
	}

	mUser.ID = id
	return mUser, nil
}
