package controller

import (
	"context"

	pb "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/http/internal"

	"google.golang.org/grpc"
)

type TournamentController interface {
	CreateUser(ctx context.Context, user *internal.User) (*internal.User, error)
	GetUserByID(ctx context.Context, id string) (*internal.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateBalanceBySum(ctx context.Context, id string, d float64) error
	LogIn(ctx context.Context, login, password string) (string, error)

	CreateTournament(ctx context.Context, name string, deposit float64) (string, error)
	GetTournamentByID(ctx context.Context, id string) (*internal.Tournament, error)
	JoinTournament(ctx context.Context, tournamentID, userID string) error
	FinishTournament(ctx context.Context, id string) error
	CancelTournament(ctx context.Context, id string) error
}

type tournamentInteractor struct {
	tgrpc pb.TournamentServiceClient
}

func NewTournamentController(cc grpc.ClientConnInterface) TournamentController {
	return &tournamentInteractor{
		tgrpc: pb.NewTournamentServiceClient(cc),
	}
}
