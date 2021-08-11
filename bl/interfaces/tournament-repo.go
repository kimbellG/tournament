package interfaces

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/debugutil"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentRepository struct{}

func (tr *TournamentRepository) Insert(repo tx.DBTX, tournament *models.Tournament) (uuid.UUID, error) {
	const query = `
		INSERT INTO Tournaments(name, deposit) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	stmt, err := repo.Prepare(query)
	if err != nil {
		return id, kerror.Newf(kerror.SQLPrepareStatementError, "prepare statement %v: %w", query, err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRow(tournament.Name, tournament.Deposit).Scan(&id); err != nil {
		return id, kerror.Newf(kerror.SQLConstraintError, "insert tournament: %w", err)
	}

	return id, nil
}

func (tr *TournamentRepository) SelectByID(repo tx.DBTX, id uuid.UUID) (*models.Tournament, error) {
	const query = `
		SELECT * FROM Tournaments WHERE id = $1
	`
	tournament := &models.Tournament{}

	stmt, err := repo.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt %v: %v", query, err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRow(id).Scan(&tournament.ID, &tournament.Name, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.Newf(kerror.TournamentDoesntExists, "tournament with id(%v) isn't exists: %v", id, err)
		}

		return nil, kerror.Newf(kerror.SQLScanError, "scan query: %v", err)
	}

	users, err := tr.selectUserIDsOfTournament(repo, id)
	if err != nil {
		return nil, kerror.Errorf(err, "get users of tournament")
	}
	tournament.Users = users

	return tournament, nil

}

func (tr *TournamentRepository) selectUserIDsOfTournament(repo tx.DBTX, tournamentID uuid.UUID) ([]models.User, error) {
	const query = `
		SELECT * FROM UsersOfTournaments WHERE tournamentID = $1;
	`
	users := []models.User{}

	stmt, err := repo.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare query: %v", err)
	}
	defer debugutil.Close(stmt)

	rows, err := stmt.Query(tournamentID)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLConstraintError, "query request: %v", err)
	}
	defer debugutil.Close(rows)

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Name, &user.Balance); err != nil {
			return nil, kerror.Newf(kerror.SQLScanError, "scan user ids of tournament(%v): %v", tournamentID, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (tr *TournamentRepository) SelectRandomUserOfTournament(repo tx.DBTX, tournamentID uuid.UUID) (*models.User, error) {
	const query = `
		WITH random_id AS (
			SELECT user FROM UsersOfTournaments WHERE tournament = $1
				OFFSET random() * COUNT(*) LIMIT 1
		)
		SELECT * FROM Users WHERER id = (SELECT id FROM random_id);
	`
	var user models.User

	stmt, err := repo.Prepare(query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRow(tournamentID).Scan(&user.ID, &user.Name, &user.Balance); err != nil {
		return nil, kerror.Newf(kerror.SQLScanError, "scan user from db: %v", err)
	}

	return &user, nil
}

func (tr *TournamentRepository) InsertUserToTournament(repo tx.DBTX, tournamentID, userID uuid.UUID) error {
	const query = `
		INSERT INTO UsersOfTournaments(tournament, user) VALUES ($1, $2); 
	`

	stmt, err := repo.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare query: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.Exec(tournamentID, userID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec stmt: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) AddToPrize(repo tx.DBTX, ID uuid.UUID, addend float64) error {
	const query = `
		UPDATE Tournaments
			SET prize = prize + $1
			WHERE id = $2
	`

	stmt, err := repo.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.Exec(addend, ID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec query: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) RefundDepositToUsers(repo tx.DBTX, tournamentID uuid.UUID) error {
	const query = `
		WITH depositOfTournament AS (
			SELECT deposit FROM Tournament WHERE id = $1
		)
		UPDATE Users
			SET balance = balance + (SELECT deposit FROM depositOfTournament)
			WHERE id IN (SELECT id FROM UsersOfTournaments WHERE tournament = $1);
	`

	stmt, err := repo.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.Exec(tournamentID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec query: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) SetWinner(repo tx.DBTX, tournamentID, winnerID uuid.UUID) error {
	const query = `
		UPDATE Tournament SET winner = $1 WHERE id = $2;
	`

	stmt, err := repo.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.Exec(winnerID, tournamentID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec update query: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) UpdateStatus(repo tx.DBTX, tournamentID uuid.UUID, newStatus models.TournamentStatus) error {
	const query = `
		UPDATE Tournament SET status = $1 WHERE id = $2;
	`

	stmt, err := repo.Prepare(query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.Exec(newStatus, tournamentID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec update query: %v", err)
	}

	return nil

}
