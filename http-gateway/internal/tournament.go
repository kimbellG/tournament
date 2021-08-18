package internal

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
	Users   []User           `json:"users"`
	Winner  string           `json:"winner"`
	Status  TournamentStatus `json:"status"`
}
