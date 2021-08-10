package interfaces

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/debugutil"
	"github.com/kimbellG/tournament-bl/models"
)

type TournamentRepository struct {
	db *sql.DB
}

func (ur *TournamentRepository) Create(tournament *models.Tournament) (uuid.UUID, error) {
	const query = `
		INSERT INTO Tournaments(name, deposit) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	stmt, err := ur.db.Prepare(query)
	if err != nil {
		return id, kerror.Newf(kerror.IntervalServerError, "prepare stmt %v: %v", query, err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRow(tournament.Name, tournament.Deposit).Scan(&id); err != nil {
		return id, kerror.Newf(kerror.BadRequest, "insert tournament: %v", err)
	}

	return id, nil
}

func (ur *TournamentRepository) GetByID(id uuid.UUID) (*models.Tournament, error) {
	const query = `
		SELECT * FROM Tournaments WHERE id = $1
	`
	tournament := &models.Tournament{}

	stmt, err := ur.db.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.IntervalServerError, "prepare stmt %v: %v", query, err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRow(id).Scan(&tournament.ID, &tournament.Name, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.Newf(kerror.InvalidID, "tournament with id(%v) isn't exists: %v", id, err)
		}

		return nil, kerror.Newf(kerror.IntervalServerError, "scan query: %v", err)
	}

	users, err := ur.getUserIDsOfTournament(id)
	if err != nil {
		return nil, kerror.Errorf(err, "get users of tournament")
	}
	tournament.Users = users

	return tournament, nil

}

func (ur *TournamentRepository) getUserIDsOfTournament(tournamentID uuid.UUID) ([]models.User, error) {
	const query = `
		SELECT * FROM UsersOfTournaments WHERE tournamentID = $1;
	`
	users := []models.User{}

	stmt, err := ur.db.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.IntervalServerError, "prepare query: %v", err)
	}
	defer debugutil.Close(stmt)

	rows, err := stmt.Query(tournamentID)
	if err != nil {
		return nil, kerror.Newf(kerror.BadRequest, "query request: %v", err)
	}
	defer debugutil.Close(rows)

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Name, &user.Balance); err != nil {
			return nil, kerror.Newf(kerror.IntervalServerError, "scan user ids of tournament(%v): %v", tournamentID, err)
		}

		users = append(users, user)
	}

	return users, nil
}
