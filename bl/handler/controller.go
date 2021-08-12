package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/controller"
	ttgrpc "github.com/kimbellG/tournament-bl/handler/grpc"
	"github.com/kimbellG/tournament-bl/models"
)

type ServiceController struct {
	userController controller.UserController
}

func (sc *ServiceController) SaveUser(ctx context.Context, user *ttgrpc.User) (*ttgrpc.SaveResponse, error) {
	mUser, err := userFromProto(user)
	if err != nil {
		return nil, kerror.Errorf(err, "marshaling user struct to models")
	}

	id, err := sc.userController.Save(mUser)
	if err != nil {
		return nil, kerror.Errorf(err, "save user")
	}

	return &ttgrpc.SaveResponse{
		Id: id.String(),
	}, nil
}

func userFromProto(gUser *ttgrpc.User) (*models.User, error) {
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

func (sc *ServiceController) GetUserById(r *ttgrpc.UserRequest) (*ttgrpc.User, error) {
	id, err := userIDFromProto(r)
	if err != nil {
		return nil, kerror.Errorf(err, "marshaling id from request")
	}

	user, err := sc.userController.GetByID(id)
	if err != nil {
		return nil, kerror.Errorf(err, "get user from controller")
	}

	return userToProto(user), nil
}

func userIDFromProto(r *ttgrpc.UserRequest) (uuid.UUID, error) {
	id, err := uuid.Parse(r.GetID())
	if err != nil {
		return id, kerror.Newf(kerror.InvalidID, "parse id: %w", err)
	}

	return id, nil
}

func userToProto(user *models.User) *ttgrpc.User {
	return &ttgrpc.User{
		ID:      user.ID.String(),
		Name:    user.Name,
		Balance: user.Balance,
	}
}
