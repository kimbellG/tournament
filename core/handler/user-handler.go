package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	ttgrpc "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/handler/kegrpc"
	"github.com/kimbellG/tournament/core/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (sc *ServiceHandler) SaveUser(ctx context.Context, user *ttgrpc.User) (*ttgrpc.SaveResponse, error) {
	savelog := log.WithFields(log.Fields{
		"action":  "save user",
		"request": user,
	})
	savelog.Info("received save user request")

	mUser := userFromProto(user)

	id, err := sc.userController.Save(ctx, mUser)
	if err != nil {
		kerror.ErrorLog(savelog, err, "Failed to save user")
		return nil, kegrpc.Errorf(err, "save user")
	}
	savelog.WithField("response", id).Info("user saved")

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
	getlog := log.WithFields(log.Fields{
		"action":  "get user by id",
		"request": r,
	})
	getlog.Info("received get user by if request")

	id, err := userIDFromProto(r)
	if err != nil {
		kerror.ErrorLog(getlog, err, "Failed to get user id from request")
		return nil, kegrpc.Errorf(err, "marshaling id from request")
	}

	user, err := sc.userController.GetByID(ctx, id)
	if err != nil {
		kerror.ErrorLog(getlog, err, "Failed to get user by id")
		return nil, kegrpc.Errorf(err, "get user from controller")
	}
	getlog.WithField("response", user).Info("get user request is successful")

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
	deletelog := log.WithFields(log.Fields{
		"action":  "delete user by id",
		"request": r,
	})
	deletelog.Info("received delete user request")

	id, err := userIDFromProto(r)
	if err != nil {
		kerror.ErrorLog(deletelog, err, "Failed to get user id from request")
		return &emptypb.Empty{}, kegrpc.Newf(kerror.InvalidID, "marshaling from user request: %w", err)
	}

	if err := sc.userController.DeleteByID(ctx, id); err != nil {
		kerror.ErrorLog(deletelog, err, "Failed to delete user")
		return &emptypb.Empty{}, kegrpc.Errorf(err, "delete user from controller")
	}
	deletelog.Info("user deleted")

	return &emptypb.Empty{}, nil
}

func (sc *ServiceHandler) SumToBalance(ctx context.Context, r *ttgrpc.RequestToUpdateBalance) (*emptypb.Empty, error) {
	sumlog := log.WithFields(log.Fields{
		"action":  "sum to user balance",
		"request": r,
	})
	sumlog.Info("received sum request")

	id, err := uuid.Parse(r.GetID())
	if err != nil {
		kerror.ErrorLog(sumlog, err, "Failed to get user id from request")
		return &emptypb.Empty{}, kegrpc.Newf(kerror.InvalidID, "parsing id from request: %w", err)
	}

	if err := sc.userController.UpdateBalance(ctx, id, r.GetAddend()); err != nil {
		kerror.ErrorLog(sumlog, err, "Failed to update user balance")
		return &emptypb.Empty{}, kegrpc.Errorf(err, "controller")
	}
	sumlog.Info("balance updated")

	return &emptypb.Empty{}, nil
}
