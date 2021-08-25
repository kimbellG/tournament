package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `sql:", type:uuid"`
	Name     string
	Password string
	Balance  float64
}
