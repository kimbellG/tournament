package internal

import "github.com/kimbellG/kerror"

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Password string  `json:"password"`
}

func (u *User) Valid() error {
	if u.Balance < 0 {
		return kerror.Newf(kerror.BadRequest, "user's balance should be more than 0")
	}

	return nil
}
