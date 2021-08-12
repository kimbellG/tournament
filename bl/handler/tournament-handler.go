package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	ttgrpc "github.com/kimbellG/tournament-bl/handler/grpc"
	"github.com/kimbellG/tournament-bl/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (sh *ServiceHandler) CreateTournament(ctx context.Context, r *ttgrpc.CreateTournamentRequest) (*ttgrpc.CreateTournamentResponse, error) {
	id, err := sh.tournamentController.Create(tournamentFromProto(r))
	if err != nil {
		return &ttgrpc.CreateTournamentResponse{}, kerror.Errorf(err, "controller")
	}

	return &ttgrpc.CreateTournamentResponse{
		Id: id.String(),
	}, nil
}

func tournamentFromProto(protoTournament *ttgrpc.CreateTournamentRequest) *models.Tournament {
	return &models.Tournament{
		Name:    protoTournament.GetName(),
		Deposit: protoTournament.GetDeposit(),
	}
}

func (sh *ServiceHandler) GetTournamentByID(ctx context.Context, r *ttgrpc.TournamentRequest) (*ttgrpc.Tournament, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return nil, kerror.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	tournament, err := sh.tournamentController.GetByID(id)
	if err != nil {
		return nil, kerror.Errorf(err, "controller")
	}

	return tournamentToProto(tournament), nil

}

func tournamentToProto(tournament *models.Tournament) *ttgrpc.Tournament {
	return &ttgrpc.Tournament{
		Id:      tournament.ID.String(),
		Name:    tournament.Name,
		Deposit: tournament.Deposit,
		Prize:   tournament.Prize,
		Users:   uuidOfUsersToStringSlice(tournament.Users),
		Winner:  tournament.Winner.String(),
		Status:  string(tournament.Status),
	}
}

func uuidOfUsersToStringSlice(users []models.User) []string {
	var uuidStrings []string
	for _, user := range users {
		uuidStrings = append(uuidStrings, user.ID.String())
	}

	return uuidStrings
}

func (sh *ServiceHandler) JoinTournament(ctx context.Context, r *ttgrpc.JoinRequest) (*emptypb.Empty, error) {
	tournament, err := uuid.Parse(r.GetTournamentID())
	if err != nil {
		return nil, kerror.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	user, err := uuid.Parse(r.GetUserID())
	if err != nil {
		return nil, kerror.Newf(kerror.InvalidID, "parsing user id: %w", err)
	}

	if err := sh.tournamentController.Join(tournament, user); err != nil {
		return nil, kerror.Errorf(err, "controller")
	}

	return nil, nil
}

func (sh *ServiceHandler) FinishTournament(ctx context.Context, r *ttgrpc.TournamentRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return nil, kerror.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	if err := sh.tournamentController.Finish(id); err != nil {
		return nil, kerror.Errorf(err, "controller")
	}

	return nil, nil
}

func (sh *ServiceHandler) CancelTournament(ctx context.Context, r *ttgrpc.TournamentRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return nil, kerror.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	if err := sh.tournamentController.Cancel(id); err != nil {
		return nil, kerror.Errorf(err, "controller")
	}

	return nil, nil
}
