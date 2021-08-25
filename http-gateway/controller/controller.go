package controller

import (
	"context"

	"github.com/kimbellG/tournament/http/internal"
)

type TournamentController interface {
	CreateUser(ctx context.Context, user *internal.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*internal.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateBalanceBySum(ctx context.Context, id string, d float64) error

	CreateTournament(ctx context.Context, name string, deposit float64) (string, error)
	GetTournamentByID(ctx context.Context, id string) (*internal.Tournament, error)
	JoinTournament(ctx context.Context, tournamentID, userID string) error
	FinishTournament(ctx context.Context, id string) error
	CancelTournament(ctx context.Context, id string) error
}
