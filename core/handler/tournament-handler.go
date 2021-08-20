package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/kerror/kegrpc"
	ttgrpc "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (sh *ServiceHandler) CreateTournament(ctx context.Context, r *ttgrpc.CreateTournamentRequest) (*ttgrpc.CreateTournamentResponse, error) {
	createLog := log.WithFields(log.Fields{
		"action":  "create tournament",
		"request": r,
	})
	createLog.Info("received create tournament request")

	id, err := sh.tournamentController.Create(ctx, tournamentFromProto(r))
	if err != nil {
		kerror.ErrorLog(createLog, err, "failed to create tournament")
		return &ttgrpc.CreateTournamentResponse{}, kegrpc.Errorf(err, "controller")
	}

	createLog.WithField("response", id).Info("create tournament request was successful")

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
	getlog := log.WithFields(log.Fields{
		"action":  "get tournement by id",
		"request": r,
	})
	getlog.Info("received get touranment by id request")

	id, err := uuid.Parse(r.GetId())
	if err != nil {
		kerror.ErrorLog(getlog, err, "Failed to parsing id of tournament")
		return nil, kegrpc.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	tournament, err := sh.tournamentController.GetByID(ctx, id)
	if err != nil {
		kerror.ErrorLog(getlog, err, "Failed request to getting tournament")
		return nil, kegrpc.Errorf(err, "controller")
	}
	getlog.WithField("response", tournament).Info("request for getting tournament was successful")

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
	joinlog := log.WithFields(log.Fields{
		"action":  "join to tournament",
		"request": r,
	})
	joinlog.Info("received join touranament request")

	tournament, err := uuid.Parse(r.GetTournamentID())
	if err != nil {
		kerror.ErrorLog(joinlog, err, "Failed to parsing id of tournament")
		return nil, kegrpc.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	user, err := uuid.Parse(r.GetUserID())
	if err != nil {
		kerror.ErrorLog(joinlog, err, "Failed to parsing joiner's id")
		return nil, kegrpc.Newf(kerror.InvalidID, "parsing user id: %w", err)
	}

	if err := sh.tournamentController.Join(ctx, tournament, user); err != nil {
		kerror.ErrorLog(joinlog, err, "Failed to join user to tournament")
		return nil, kegrpc.Errorf(err, "controller")
	}

	joinlog.Info("user joined to tournament")

	return &emptypb.Empty{}, nil
}

func (sh *ServiceHandler) FinishTournament(ctx context.Context, r *ttgrpc.TournamentRequest) (*emptypb.Empty, error) {
	finishlog := log.WithFields(log.Fields{
		"action":  "finish tournament",
		"request": r,
	})
	finishlog.Info("received finish tournament request")

	id, err := uuid.Parse(r.GetId())
	if err != nil {
		kerror.ErrorLog(finishlog, err, "Failed to parsing tournament id")
		return nil, kegrpc.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	if err := sh.tournamentController.Finish(ctx, id); err != nil {
		kerror.ErrorLog(finishlog, err, "Failed to finish tournament")
		return nil, kegrpc.Errorf(err, "controller")
	}
	finishlog.Info("tournament finished")

	return &emptypb.Empty{}, nil
}

func (sh *ServiceHandler) CancelTournament(ctx context.Context, r *ttgrpc.TournamentRequest) (*emptypb.Empty, error) {
	cancellog := log.WithFields(log.Fields{
		"action":  "cancel tournament",
		"request": r,
	})
	cancellog.Info("received cancel tournament request")

	id, err := uuid.Parse(r.GetId())
	if err != nil {
		kerror.ErrorLog(cancellog, err, "Failed to parsing tournament id")
		return nil, kegrpc.Newf(kerror.InvalidID, "parsing tournament id: %w", err)
	}

	if err := sh.tournamentController.Cancel(ctx, id); err != nil {
		kerror.ErrorLog(cancellog, err, "Failed to cancel tournament")
		return nil, kegrpc.Errorf(err, "controller")
	}

	cancellog.Info("tournament canceled")

	return &emptypb.Empty{}, nil
}
