package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/kerror/kegrpc"
	ttgrpc "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (sc *ServiceHandler) SaveUser(ctx context.Context, user *ttgrpc.User) (*ttgrpc.SaveResponse, error) {
	mUser := userFromProto(user)

	id, err := sc.userController.Save(ctx, mUser)
	if err != nil {
		return nil, kegrpc.Errorf(err, "save user")
	}

	return &ttgrpc.SaveResponse{
		Id: id.String(),
	}, nil
}

func userFromProto(gUser *ttgrpc.User) *models.User {
	return &models.User{
		Name:    gUser.GetName(),
		Balance: gUser.GetBalance(),
	}
}

func (sc *ServiceHandler) GetUserById(ctx context.Context, r *ttgrpc.UserRequest) (*ttgrpc.User, error) {
	id, err := userIDFromProto(r)
	if err != nil {
		return nil, kegrpc.Errorf(err, "marshaling id from request")
	}

	user, err := sc.userController.GetByID(ctx, id)
	if err != nil {
		return nil, kegrpc.Errorf(err, "get user from controller")
	}

	return userToProto(user), nil
}

func userIDFromProto(r *ttgrpc.UserRequest) (uuid.UUID, error) {
	id, err := uuid.Parse(r.GetID())
	if err != nil {
		return id, kegrpc.Newf(kerror.InvalidID, "parse id: %w", err)
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

func (sc *ServiceHandler) DeleteUserByID(ctx context.Context, r *ttgrpc.UserRequest) (*emptypb.Empty, error) {
	id, err := userIDFromProto(r)
	if err != nil {
		return &emptypb.Empty{}, kegrpc.Newf(kerror.InvalidID, "marshaling from user request: %w", err)
	}

	if err := sc.userController.DeleteByID(ctx, id); err != nil {
		return &emptypb.Empty{}, kegrpc.Errorf(err, "delete user from controller")
	}

	return &emptypb.Empty{}, nil
}

func (sc *ServiceHandler) SumToBalance(ctx context.Context, r *ttgrpc.RequestToUpdateBalance) (*emptypb.Empty, error) {
	id, err := uuid.Parse(r.GetID())
	if err != nil {
		return &emptypb.Empty{}, kegrpc.Newf(kerror.InvalidID, "parsing id from request: %w", err)
	}

	if err := sc.userController.UpdateBalance(ctx, id, r.GetAddend()); err != nil {
		return &emptypb.Empty{}, kegrpc.Errorf(err, "controller")
	}

	return &emptypb.Empty{}, nil
}
