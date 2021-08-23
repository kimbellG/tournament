package models

import (
	"github.com/google/uuid"
)

type TournamentStatus string

const (
	Active TournamentStatus = "Active"
	Cancel TournamentStatus = "Cancel"
	Finish TournamentStatus = "Finish"
)

type Tournament struct {
	ID      uuid.UUID `sql:", type:uuid"`
	Name    string
	Deposit float64
	Prize   float64
	Users   []User
	Winner  uuid.UUID
	Status  TournamentStatus
}
