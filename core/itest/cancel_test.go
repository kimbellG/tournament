// +build integration

package itest

import (
	"context"
	"fmt"
	"testing"

	tgrpc "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestCancelTournament(t *testing.T) {
	conn, db := initTest(t)
	defer conn.Close()
	defer db.Close()

	client := tgrpc.NewTournamentServiceClient(conn)

	activeTournament := createTournament(t, db, &models.Tournament{
		Name:    "tournament to cancel",
		Deposit: 1000,
		Prize:   10000,
		Status:  models.Active,
	})

	notActiveTournament := createTournament(t, db, &models.Tournament{
		Name:    "not active tournament",
		Deposit: 500,
		Prize:   1000,
		Status:  models.Cancel,
	})

	var users []*models.User
	for i := 0; i < 4; i++ {
		user := createUser(t, db, &models.User{
			Name:    fmt.Sprintf("cancel user %d", i),
			Balance: 500,
		})
		users = append(users, user)

		if _, err := db.Exec("INSERT INTO UsersOfTournaments(tournamentID, userID) VALUES($1, $2)", activeTournament.ID, user.ID); err != nil {
			t.Fatalf("Failed to join user to tournament: %v", err)
		}

		if _, err := db.Exec("INSERT INTO UsersOfTournaments(tournamentID, userID) VALUES($1, $2)", notActiveTournament.ID, user.ID); err != nil {
			t.Fatalf("Failed to join user to tournament: %v", err)
		}
	}

	tt := []struct {
		name       string
		tournament *models.Tournament
		users      []*models.User
		code       codes.Code
	}{
		{
			name:       "active tournament to cancel",
			tournament: activeTournament,
			users:      users,
			code:       codes.OK,
		},
		{
			name:       "not active tournament to cancel",
			tournament: notActiveTournament,
			users:      users,
			code:       codes.InvalidArgument,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := client.CancelTournament(context.Background(), &tgrpc.TournamentRequest{Id: tc.tournament.ID.String()}); err != nil {
				assertGrpcError(t, tc.code, err)
				return
			}

			for _, user := range tc.users {
				var actualBalance float64
				if err := db.QueryRow("SELECT balance FROM Users WHERE id = $1", user.ID).Scan(&actualBalance); err != nil {
					t.Fatalf("Failed to select joiner's balance. joiner: %v; tournament: %v", user.Name, tc.tournament.Name)
				}

				assert.Equalf(t, user.Balance+tc.tournament.Deposit, actualBalance, "actual balance should be old balance(%v) + tournament deposit(%v)", user.Balance, tc.tournament.Deposit)
			}

			var status models.TournamentStatus
			if err := db.QueryRow("SELECT status FROM Tournaments WHERE id = $1", tc.tournament.ID).Scan(&status); err != nil {
				t.Fatalf("Failed to select status of tournament(%v): %v", tc.tournament.Status, err)
			}

			assert.Equal(t, models.Cancel, status, "tournament status should be cancel")
		})
	}
}
