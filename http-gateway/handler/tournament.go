package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kimbellG/tournament/http/internal"
)

type CreateTournamentResponse struct {
	ID string `json:"id"`
}

func (h *Handler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	tournament := &internal.Tournament{}
	if err := json.NewDecoder(r.Body).Decode(tournament); err != nil {
		http.Error(w, "Failed to decode create's request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.tournament.CreateTournament(r.Context(), tournament.Name, tournament.Deposit)
	if err != nil {
		http.Error(w, "Failed to create tournament: "+err.Error(), parsing)
		return
	}

	if err := json.NewEncoder(w).Encode(CreateTournamentResponse{id}); err != nil {
		http.Error(w, "Failed to encode answer to response body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTournamentByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[idPath]

	tournament, err := h.tournament.GetTournamentByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get tournament by id: "+err.Error(), parsing)
		return
	}

	if err := json.NewEncoder(w).Encode(tournament); err != nil {
		http.Error(w, "Failed to encode tournament in response body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

type JoinRequest struct {
	UserID string `json:"userId"`
}

func (h *Handler) JoinTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := mux.Vars(r)[idPath]
	joinRequest := &JoinRequest{}

	if err := json.NewDecoder(r.Body).Decode(joinRequest); err != nil {
		http.Error(w, "Failed to decode join request body:"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.tournament.JoinTournament(r.Context(), tournamentID, joinRequest.UserID); err != nil {
		http.Error(w, "Failed to join user to tournament: "+err.Error(), parsing)
		return
	}
}

func (h *Handler) FinishTournament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[idPath]

	if err := h.tournament.FinishTournament(r.Context(), id); err != nil {
		http.Error(w, "Failed to finish tournament: "+err.Error(), parsing)
		return
	}
}

func (h *Handler) CancelTournament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[idPath]

	if err := h.tournament.CancelTournament(r.Context(), id); err != nil {
		http.Error(w, "Failed to cancel tournament: "+err.Error(), parsing)
		return
	}
}
