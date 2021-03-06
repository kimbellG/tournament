package controller

import (
	"context"

	"github.com/kimbellG/kerror"
	pb "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/http/internal"
)

func (t *tournamentInteractor) CreateUser(ctx context.Context, user *internal.User) (*internal.User, error) {
	resp, err := t.tgrpc.SaveUser(ctx, userToProto(user))
	if err != nil {
		return nil, kerror.Errorf(err, "grpc-core")
	}
	user.ID = resp.GetId()
	user.Password = resp.GetPassword()

	return user, nil
}

func userToProto(user *internal.User) *pb.User {
	return &pb.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}
}

func (t *tournamentInteractor) GetUserByID(ctx context.Context, id string) (*internal.User, error) {
	resp, err := t.tgrpc.GetUserByID(ctx, &pb.UserRequest{ID: id})
	if err != nil {
		return nil, kerror.Errorf(err, "grpc-core")
	}

	return userFromProto(resp), nil
}

func userFromProto(user *pb.User) *internal.User {
	return &internal.User{
		ID:      user.GetID(),
		Name:    user.GetName(),
		Balance: user.GetBalance(),
	}
}

func (t *tournamentInteractor) DeleteUser(ctx context.Context, id string) error {
	if _, err := t.tgrpc.DeleteUserByID(ctx, &pb.UserRequest{ID: id}); err != nil {
		return kerror.Errorf(err, "grpc-core")
	}

	return nil
}

func (t *tournamentInteractor) UpdateBalanceBySum(ctx context.Context, id string, d float64) error {
	if _, err := t.tgrpc.SumToBalance(ctx, &pb.RequestToUpdateBalance{ID: id, Addend: d}); err != nil {
		return kerror.Errorf(err, "grpc-core")
	}

	return nil
}

func (t *tournamentInteractor) LogIn(ctx context.Context, login, password string) (string, error) {
	resp, err := t.tgrpc.UserAuthorization(ctx, &pb.AuthorizationRequest{Username: login, Password: password})
	if err != nil {
		return "", kerror.Errorf(err, "grpc request to core service")
	}

	return resp.GetId(), nil
}
