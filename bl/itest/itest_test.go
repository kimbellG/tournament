// +build integration

package itest

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	tgrpc "github.com/kimbellG/tournament-bl/handler/grpc"
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

func TestCreateTournament(t *testing.T) {
	conn, db := initTest(t)
	defer conn.Close()
	defer db.Close()

	client := tgrpc.NewTournamentServiceClient(conn)

	tt := []struct {
		name    string
		input   tgrpc.CreateTournamentRequest
		errCode codes.Code
	}{
		{
			"successful tournament", tgrpc.CreateTournamentRequest{Name: "artyom tournament", Deposit: 1000}, codes.OK,
		},
		{
			"negative deposit", tgrpc.CreateTournamentRequest{Name: "failed tournament", Deposit: -100}, codes.FailedPrecondition,
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

func assertGrpcError(t *testing.T, wantCode codes.Code, err error) {
	if e, ok := status.FromError(err); ok {
		assert.Equalf(t, wantCode, e.Code(), "Unexcepted error with %s status code: %v", e.Code().String(), e.Err())
	} else {
		t.Errorf("Unknown error without status code: %v", err)
	}
}

func TestFinishTournament(t *testing.T) {
	conn, db := initTest(t)
	defer conn.Close()
	defer db.Close()

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
			name:       "Active touranment",
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

func isUserInTournament(userID uuid.UUID, users []*models.User) (*models.User, bool) {
	for _, user := range users {
		if user.ID == userID {
			return user, true
		}
	}

	return nil, false
}
