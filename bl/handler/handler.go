package handler

import (
	"github.com/kimbellG/tournament-bl/controller"
)

type ServiceHandler struct {
	userController       controller.UserController
	tournamentController controller.TournamentController
}

func NewServiceHandler(user controller.UserController, tournament controller.TournamentController) *ServiceHandler {
	return &ServiceHandler{
		userController:       user,
		tournamentController: tournament,
	}
}
