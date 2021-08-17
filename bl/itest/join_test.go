// +build integration

package itest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	tgrpc "github.com/kimbellG/tournament-bl/handler/grpc"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestJoinToTournament(t *testing.T) {
	conn, db := initTest(t)
	defer conn.Close()
	defer db.Close()

	client := tgrpc.NewTournamentServiceClient(conn)

	activeTournament := createTournament(t, db, &models.Tournament{
		Name:    "join tournament",
		Deposit: 1000,
		Status:  models.Active,
	})
	notActiveTournament := createTournament(t, db, &models.Tournament{
		Name:    "not active tournament",
		Deposit: 10,
		Status:  models.Finish,
	})

	okUser := createUser(t, db, &models.User{
		Name:    "joinOK",
		Balance: 1500,
	})
	withoutDepositUser := createUser(t, db, &models.User{
		Name:    "without deposit user",
		Balance: 500,
	})

	tt := []struct {
		name           string
		wantTournament *models.Tournament
		wantUser       *models.User
		wantPrize      float64
		code           codes.Code
	}{
		{
			name:           "Success",
			wantTournament: activeTournament,
			wantUser:       okUser,
			wantPrize:      1000,
			code:           codes.OK,
		},
		{
			name:           "user without deposit",
			wantTournament: activeTournament,
			wantUser:       withoutDepositUser,
			wantPrize:      1000,
			code:           codes.FailedPrecondition,
		},
		{
			name:           "not active tournament",
			wantTournament: notActiveTournament,
			wantUser:       okUser,
			code:           codes.InvalidArgument,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.JoinTournament(context.Background(), &tgrpc.JoinRequest{TournamentID: tc.wantTournament.ID.String(), UserID: tc.wantUser.ID.String()})
			if err != nil {
				assertGrpcError(t, tc.code, err)
				return
			}

			var actualBalance float64
			if err := db.QueryRow("SELECT balance FROM Users WHERE id = $1", tc.wantUser.ID).Scan(&actualBalance); err != nil {
				t.Fatalf("Failed select balance of the want user from database: %v", err)
			}
			assert.Equalf(t, tc.wantUser.Balance-tc.wantTournament.Deposit, actualBalance,
				"New balance should be old balance(%v) - deposit of tournament(%v)", tc.wantUser.Balance, tc.wantTournament.Deposit,
			)

			var actualPrize float64
			if err := db.QueryRow("SELECT prize FROM Tournaments WHERE id = $1", tc.wantTournament.ID).Scan(&actualPrize); err != nil {
				t.Fatalf("Failed select prize from tournament: %v", err)
			}
			assert.Equalf(t, tc.wantPrize, actualPrize,
				"New prize should be old prize + deposit(%v)",
				tc.wantTournament.Deposit,
			)

			var joinerID uuid.UUID
			if err := db.QueryRow("SELECT id FROM UsersOfTournaments WHERE tournamentID = $1 AND userID = $2", tc.wantTournament.ID, tc.wantUser.ID).Scan(&joinerID); err != nil {
				if !assert.Equalf(t, err, sql.ErrNoRows, "user(%v) in tournament(%v) not found", tc.wantTournament.ID, tc.wantUser.ID) {
					t.Errorf("failed to select joiner from database: %v", err)
				}
			}
		})
	}
}
