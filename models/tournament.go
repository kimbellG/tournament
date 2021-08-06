package models

type Tournament struct {
	ID      string
	Name    string
	Deposit float64
	Prize   float64
	Users   []Users
	Winner  string
}
