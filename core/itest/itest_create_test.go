// +build integration

package itest

import (
	"context"
	"database/sql"
	"testing"

	tgrpc "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestCreateTournament(t *testing.T) {
	client := tgrpc.NewTournamentServiceClient(conn)

	tt := []struct {
		name    string
		input   tgrpc.CreateTournamentRequest
		errCode codes.Code
	}{
		{
			"successful tournament", tgrpc.CreateTournamentRequest{
				Name:    "artyom tournament",
				Deposit: 1000,
			},
			codes.OK,
		},
		{
			"negative deposit",
			tgrpc.CreateTournamentRequest{
				Name:    "failed tournament",
				Deposit: -100,
			},
			codes.FailedPrecondition,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id, err := client.CreateTournament(context.Background(), &tc.input)
			if err != nil {
				assertGrpcError(t, tc.errCode, err)
				return
			}

			row := db.QueryRow("SELECT name, deposit FROM Tournaments WHERE id = $1", id.GetId())

			var result models.Tournament
			if err := row.Scan(&result.Name, &result.Deposit); err != nil {
				if assert.NotEqualf(t, err, sql.ErrNoRows, "Create failed. result tournament not found: Name=%v; Deposit=%v", result.Name, result.Deposit) {
					t.Errorf("Error with database: scan object: %v", err)
				}
			}

		})
	}
}
