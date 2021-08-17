// +build integration

package itest

import (
	"context"
	"fmt"
	"testing"

	tgrpc "github.com/kimbellG/tournament-bl/handler/grpc"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestGetTournamentByID(t *testing.T) {
	conn, db := initTest(t)
	defer conn.Close()
	defer db.Close()

	client := tgrpc.NewTournamentServiceClient(conn)

	tournament := createTournament(t, db, &models.Tournament{
		Name:    "ok",
		Deposit: 1203,
		Status:  models.Active,
	})

	joiners := []models.User{}
	for i := 0; i < 4; i++ {
		user := createUser(t, db, &models.User{
			Name:    fmt.Sprintf("userForGet%d", i),
			Balance: 0,
		})

		joiners = append(joiners, *user)

		if _, err := db.Exec("INSERT INTO UsersOfTournaments(tournamentID, userID) VALUES($1, $2)", tournament.ID, user.ID); err != nil {
			t.Fatalf("Failed join user to tournament in database(%d): %v", i, err)
		}
	}

	tt := []struct {
		name string
		want models.Tournament
		id   string
		code codes.Code
	}{
		{
			name: "ok tournament",
			want: models.Tournament{
				ID:      tournament.ID,
				Name:    tournament.Name,
				Deposit: tournament.Deposit,
				Users:   joiners,
			},
			id:   tournament.ID.String(),
			code: codes.OK,
		},
		{
			name: "invalid id",
			id:   "asdasd-asdasd-asda-sdasd-as",
			code: codes.InvalidArgument,
		},
		{
			name: "non-existent tournament",
			id:   "123e4567-e89b-12d3-a456-426614174000",
			code: codes.NotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tournament, err := client.GetTournamentByID(context.Background(), &tgrpc.TournamentRequest{Id: tc.id})
			if err != nil {
				assertGrpcError(t, tc.code, err)
				return
			}

			assert.Equal(t, tc.want.ID.String(), tournament.GetId(), "tournament id should be equal")
			assert.Equal(t, tc.want.Name, tournament.GetName(), "tournament name should be equal")
			compareJoiners(t, tc.want.Users, tournament.GetUsers())
		})
	}
}
