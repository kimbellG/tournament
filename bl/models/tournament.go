package models

import (
	"github.com/google/uuid"
)

type Tournament struct {
	ID      uuid.UUID `"sql:", type:uuid"`
	Name    string
	Deposit float64
	Prize   float64
	Users   []User
	Winner  string
}
