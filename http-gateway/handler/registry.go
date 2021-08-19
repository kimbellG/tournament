package handler

import (
	"github.com/gorilla/mux"
	"github.com/kimbellG/tournament/http/controller"
)

const idPath = "id"

func RegisterUserEndpoints(router *mux.Router, tc controller.TournamentController) {
	h := NewHandler(tc)

	router.HandleFunc("/", h.CreateUser).Methods("POST")
	router.HandleFunc("/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", h.GetUserByID).Methods("GET")
	router.HandleFunc("/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", h.DeleteUser).Methods("DELETE")
}
