package handler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament/http/internal"

	"github.com/gorilla/mux"
)

type CreateUserResponse struct {
	ID       string `json:"id"`
	Password string `json:"password"`
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

	created, err := h.tournament.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create user:"+err.Error(), decodeStatusCode(err))
		return
	}

	if err := json.NewEncoder(w).Encode(CreateUserResponse{ID: created.ID, Password: created.Password}); err != nil {
		http.Error(w, "Failed to encode answer in body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetUserResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)

	user, err := h.tournament.GetUserByID(r.Context(), id[IDPath])
	if err != nil {
		http.Error(w, "Failed to get user by id: "+err.Error(), decodeStatusCode(err))
		return
	}

	resp := &GetUserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode user in body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)

	if err := h.tournament.DeleteUser(r.Context(), id[IDPath]); err != nil {
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
	id := mux.Vars(r)[IDPath]

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
	id := mux.Vars(r)[IDPath]

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

type LogInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LogInResponse struct {
	Token string `json:"token"`
}

type LogClaims struct {
	ID string
	jwt.StandardClaims
}

func (h *Handler) UserLogIn(w http.ResponseWriter, r *http.Request) {
	rBody := &LogInRequest{}

	if err := json.NewDecoder(r.Body).Decode(rBody); err != nil {
		http.Error(w, "Failed to decode log in request: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.tournament.LogIn(r.Context(), rBody.Login, passwordHash(rBody.Password))
	if err != nil {
		http.Error(w, "Failed to user log in: "+err.Error(), decodeStatusCode(err))
		return
	}

	tk, err := createToken(id)
	if err != nil {
		http.Error(w, "Failed to create token: "+err.Error(), decodeStatusCode(err))
		return
	}

	if err := json.NewEncoder(w).Encode(&LogInResponse{Token: tk}); err != nil {
		http.Error(w, "Failed to encode reponse: ", http.StatusInternalServerError)
		return
	}

}

func passwordHash(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func createToken(id string) (string, error) {
	claims := &LogClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tkString, err := tk.SignedString([]byte(os.Getenv("TK_PASSWORD")))
	if err != nil {
		return "", kerror.Newf(kerror.InternalServerError, "create string from token struct: %v", err)
	}

	return tkString, nil

}
