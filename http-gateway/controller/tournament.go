package controller

import (
	"context"

	"github.com/kimbellG/kerror"
	pb "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/http/internal"
)

func (t *tournamentInteractor) CreateTournament(ctx context.Context, name string, deposit float64) (string, error) {
	resp, err := t.tgrpc.CreateTournament(ctx, &pb.CreateTournamentRequest{Name: name, Deposit: deposit})
	if err != nil {
		return "", kerror.Errorf(decodeGrpcError(err), "grcp-core")
	}

	return resp.GetId(), nil
}

func (t *tournamentInteractor) GetTournamentByID(ctx context.Context, id string) (*internal.Tournament, error) {
	tournament, err := t.tgrpc.GetTournamentByID(ctx, &pb.TournamentRequest{Id: id})
	if err != nil {
		return nil, kerror.Errorf(decodeGrpcError(err), "grpc-core")
	}

	return tournamentFromGrpc(tournament), nil
}

func tournamentFromGrpc(tournament *pb.Tournament) *internal.Tournament {
	return &internal.Tournament{
		ID:      tournament.GetId(),
		Name:    tournament.GetName(),
		Deposit: tournament.GetDeposit(),
		Prize:   tournament.GetPrize(),
		Users:   tournament.GetUsers(),
		Winner:  tournament.GetWinner(),
	}
}

func (t *tournamentInteractor) JoinTournament(ctx context.Context, tournamentID, userID string) error {
	if _, err := t.tgrpc.JoinTournament(ctx, &pb.JoinRequest{TournamentID: tournamentID, UserID: userID}); err != nil {
		return kerror.Errorf(decodeGrpcError(err), "grpc-core")
	}

	return nil
}

func (t *tournamentInteractor) FinishTournament(ctx context.Context, tournamentID string) error {
	if _, err := t.tgrpc.FinishTournament(ctx, &pb.TournamentRequest{Id: tournamentID}); err != nil {
		return kerror.Errorf(decodeGrpcError(err), "grpc-core")
	}

	return nil
}

func (t *tournamentInteractor) CancelTournament(ctx context.Context, tournamentID string) error {
	if _, err := t.tgrpc.CancelTournament(ctx, &pb.TournamentRequest{Id: tournamentID}); err != nil {
		return kerror.Errorf(decodeGrpcError(err), "grpc-core")
	}

	return nil
}
