package handler

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/kimbellG/tournament/http/controller"
)

const idPath = "id"
const uuidRegex = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func RegisterUserEndpoints(router *mux.Router, tc controller.TournamentController) {
	h := NewHandler(tc)

	router.HandleFunc("/user", h.CreateUser).Methods("POST")
	router.HandleFunc(fmt.Sprintf("/user/{id:%s}", uuidRegex), h.GetUserByID).Methods("GET")
	router.HandleFunc(fmt.Sprintf("/user/{id:%s}", uuidRegex), h.DeleteUser).Methods("DELETE")
	router.HandleFunc(fmt.Sprintf("/user/{id:%s}/take", uuidRegex), h.TakeFromBalance).Methods("POST")
	router.HandleFunc(fmt.Sprintf("/user/{id:%s}/fund", uuidRegex), h.AddToBalance).Methods("POST")
}

func RegisterTournamentEndpoints(router *mux.Router, tc controller.TournamentController) {
	h := NewHandler(tc)

	router.HandleFunc("/tournament", h.CreateTournament).Methods("POST")
	router.HandleFunc(fmt.Sprintf("/tournament/{id:%s}", uuidRegex), h.GetTournamentByID).Methods("GET")
	router.HandleFunc(fmt.Sprintf("/tournament/{id:%s}", uuidRegex), h.CancelTournament).Methods("DELETE")
	router.HandleFunc(fmt.Sprintf("/tournament/{id:%s}/join", uuidRegex), h.JoinTournament).Methods("POST")
	router.HandleFunc(fmt.Sprintf("/tournament/{id:%s}/finish", uuidRegex), h.JoinTournament).Methods("POST")

}
