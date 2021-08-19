package handler

import (
	"github.com/kimbellG/tournament/http/controller"
)

type Handler struct {
	tournament controller.TournamentController
}

func NewHandler(tournament controller.TournamentController) *Handler {
	return &Handler{
		tournament: tournament,
	}
}
