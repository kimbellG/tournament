package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament/http/internal"

	"github.com/gorilla/mux"
)

type CreateUserResponse struct {
	ID string `json:"id"`
}

func Close(cl io.Closer) {
	if err := cl.Close(); err != nil {
		log.Println(err)
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &internal.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Failed to marshal json request for creating user", http.StatusBadRequest)
		return
	}
	defer Close(r.Body)

	if err := user.Valid(); err != nil {
		http.Error(w, "Failed to validate user request: "+err.Error(), decodeStatusCode(err))
		return
	}

	id, err := h.tournament.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create user:"+err.Error(), decodeStatusCode(err))
		return
	}

	if err := json.NewEncoder(w).Encode(CreateUserResponse{id}); err != nil {
		http.Error(w, "Failed to encode answer in body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)

	user, err := h.tournament.GetUserByID(r.Context(), id[idPath])
	if err != nil {
		http.Error(w, "Failed to get user by id: "+err.Error(), decodeStatusCode(err))
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user in body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)

	if err := h.tournament.DeleteUser(r.Context(), id[idPath]); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), decodeStatusCode(err))
		return
	}
}

type UpdateBalanceRequest struct {
	Summand float64 `json:"summand"`
}

func (u *UpdateBalanceRequest) Valid() error {
	if u.Summand <= 0 {
		return kerror.Newf(kerror.BadRequest, "summand to balance should be more than 0")
	}

	return nil
}

func (h *Handler) AddToBalance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[idPath]

	updateRequest := &UpdateBalanceRequest{}
	if err := json.NewDecoder(r.Body).Decode(updateRequest); err != nil {
		http.Error(w, "Failed to decode update request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := updateRequest.Valid(); err != nil {
		http.Error(w, "Failed to validate add request: "+err.Error(), decodeStatusCode(err))
		return
	}

	if err := h.tournament.UpdateBalanceBySum(r.Context(), id, updateRequest.Summand); err != nil {
		http.Error(w, "Failed to add points to user balance: "+err.Error(), decodeStatusCode(err))
		return
	}
}

func (h *Handler) TakeFromBalance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[idPath]

	takeRequest := &UpdateBalanceRequest{}
	if err := json.NewDecoder(r.Body).Decode(takeRequest); err != nil {
		http.Error(w, "Failed to decode take request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := takeRequest.Valid(); err != nil {
		http.Error(w, "Failed to validate take request: "+err.Error(), decodeStatusCode(err))
		return
	}

	if err := h.tournament.UpdateBalanceBySum(r.Context(), id, -takeRequest.Summand); err != nil {
		http.Error(w, "Failed to take points from user balance: "+err.Error(), decodeStatusCode(err))
		return
	}
}
