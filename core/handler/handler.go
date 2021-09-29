package handler

import (
	"github.com/kimbellG/tournament/core/controller"
	ttgrpc "github.com/kimbellG/tournament/core/handler/grpc"
)

type ServiceHandler struct {
	ttgrpc.UnimplementedTournamentServiceServer

	userController       controller.UserController
	tournamentController controller.TournamentController
}

func NewServiceHandler(user controller.UserController, tournament controller.TournamentController) *ServiceHandler {
	return &ServiceHandler{
		userController:       user,
		tournamentController: tournament,
	}
}
