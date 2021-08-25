package internal

import "github.com/kimbellG/kerror"

type TournamentStatus string

const (
	Active TournamentStatus = "Active"
	Cancel TournamentStatus = "Cancel"
	Finish TournamentStatus = "Finish"
)

type Tournament struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Deposit float64          `json:"deposit"`
	Prize   float64          `json:"prize"`
	Users   []string         `json:"users"`
	Winner  string           `json:"winner"`
	Status  TournamentStatus `json:"status"`
}

func (t *Tournament) Valid() error {
	if t.Deposit <= 0 {
		return kerror.Newf(kerror.BadRequest, "deposit should be more than 0")
	}

	if t.Prize < 0 {
		return kerror.Newf(kerror.BadRequest, "prize should be positive")
	}

	return nil
}
