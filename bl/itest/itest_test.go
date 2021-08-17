// +build integration

package itest

import (
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kimbellG/tournament-bl/models"
	tournament "github.com/kimbellG/tournament-bl/service/core"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestMain(m *testing.M) {
	go tournament.StartServer()

	code := m.Run()
	os.Exit(code)
}

func initTest(t *testing.T) (*grpc.ClientConn, *sql.DB) {
	db, err := tournament.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v\n", err)
	}

	conn, err := grpc.Dial(os.Getenv("SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to initializate grpc connection: %v\n", err)
	}

	return conn, db
}

func createTournament(t *testing.T, db *sql.DB, tournament *models.Tournament) *models.Tournament {
	if err := db.QueryRow("INSERT INTO Tournaments(name, deposit, prize, status) VALUES ($1, $2, $3, $4) RETURNING id",
		tournament.Name,
		tournament.Deposit,
		tournament.Prize,
		tournament.Status).Scan(&tournament.ID); err != nil {
		t.Fatalf("Failed insert tournament in database: %v", err)
	}

	return tournament
}

func createUser(t *testing.T, db *sql.DB, user *models.User) *models.User {
	if err := db.QueryRow("INSERT INTO Users(name, balance) VALUES ($1, $2) RETURNING id", user.Name, user.Balance).Scan(&user.ID); err != nil {
		t.Fatalf("Failed insert user in database(%v): %v", user.Name, err)
	}

	return user
}

func compareJoiners(t *testing.T, excepted []models.User, actual []string) bool {
	if !assert.Equal(t, len(excepted), len(actual), "Length of joiners should be equal") {
		return false
	}

	var exceptedStrings []string
	for _, user := range excepted {
		exceptedStrings = append(exceptedStrings, user.ID.String())
	}

	return assert.Equal(t, exceptedStrings, actual, "joiners should be equal")
}

func assertGrpcError(t *testing.T, wantCode codes.Code, err error) {
	if e, ok := status.FromError(err); ok {
		assert.Equalf(t, wantCode, e.Code(), "Unexcepted error with %s status code: %v", e.Code().String(), e.Err())
	} else {
		t.Errorf("Unknown error without status code: %v", err)
	}
}

func isUserInTournament(userID uuid.UUID, users []*models.User) (*models.User, bool) {
	for _, user := range users {
		if user.ID == userID {
			return user, true
		}
	}

	return nil, false
}
