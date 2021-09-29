package handler

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/kimbellG/tournament/http/controller"
)

const (
	IDPath         = "id"
	UserPath       = "user"
	TournamentPath = "tournament"
	LogInPath      = "login"
)

const uuidRegex = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func RegisterUserEndpoints(router *mux.Router, tc controller.TournamentController) {
	h := NewHandler(tc)

	router.HandleFunc(fmt.Sprintf("/%s", UserPath),
		h.CreateUser).Methods("POST")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}", UserPath, IDPath, uuidRegex),
		h.GetUserByID).Methods("GET")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}", UserPath, IDPath, uuidRegex),
		h.DeleteUser).Methods("DELETE")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}/take", UserPath, IDPath, uuidRegex),
		h.TakeFromBalance).Methods("POST")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}/fund", UserPath, IDPath, uuidRegex),
		h.AddToBalance).Methods("POST")

	router.HandleFunc(fmt.Sprintf("/%s", LogInPath),
		h.UserLogIn).Methods("GET")
}

func RegisterTournamentEndpoints(router *mux.Router, tc controller.TournamentController) {
	h := NewHandler(tc)

	router.HandleFunc(fmt.Sprintf("/%s", TournamentPath),
		h.CreateTournament).Methods("POST")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}", TournamentPath, IDPath, uuidRegex),
		h.GetTournamentByID).Methods("GET")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}", TournamentPath, IDPath, uuidRegex),
		h.CancelTournament).Methods("DELETE")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}/join", TournamentPath, IDPath, uuidRegex),
		h.JoinTournament).Methods("POST")

	router.HandleFunc(fmt.Sprintf("/%s/{%s:%s}/finish", TournamentPath, IDPath, uuidRegex),
		h.JoinTournament).Methods("POST")
}
