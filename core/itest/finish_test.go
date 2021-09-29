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

func TestFinishTournament(t *testing.T) {
	client := tgrpc.NewTournamentServiceClient(conn)

	activeTournament := createTournament(t, db, &models.Tournament{
		Name:    "finish tournament",
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
			Name:    fmt.Sprintf("finish user %d", i),
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
			name:       "Active tournament",
			tournament: activeTournament,
			users:      users,
			code:       codes.OK,
		},
		{
			name:       "Not active tournament",
			tournament: notActiveTournament,
			users:      users,
			code:       codes.InvalidArgument,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := client.FinishTournament(context.Background(), &tgrpc.TournamentRequest{Id: tc.tournament.ID.String()}); err != nil {
				assertGrpcError(t, tc.code, err)
				return
			}

			winner := &models.User{}
			if err := db.QueryRow("SELECT users.id, users.name, users.balance FROM Tournaments INNER JOIN Users ON winner = users.id WHERE tournaments.id = $1",
				tc.tournament.ID).Scan(&winner.ID, &winner.Name, &winner.Balance); err != nil {
				t.Fatalf("Failed select winner from tournament: %v", err)
			}

			oldUser, ok := isUserInTournament(winner.ID, tc.users)
			if !assert.Truef(t, ok, "winner(%v) isn't tournament joiner(%v)", winner.Name, tc.tournament.Name) {
				return
			}

			assert.Equalf(t, oldUser.Balance+tc.tournament.Prize, winner.Balance, "new balance winner should be old balance(%v) + prize(%v)", oldUser.Balance, tc.tournament.Prize)

			var status models.TournamentStatus
			if err := db.QueryRow("SELECT status FROM Tournaments WHERE id = $1", tc.tournament.ID).Scan(&status); err != nil {
				t.Fatalf("Failed select status of tournament: %v", err)
			}
			assert.Equal(t, status, models.Finish, "status should be finish")

		})
	}
}
