package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament-bl/debugutil"
	"github.com/kimbellG/tournament-bl/models"
	"github.com/kimbellG/tournament-bl/tx"
)

type TournamentRepository struct{}

func (tr *TournamentRepository) Insert(ctx context.Context, store tx.DBTX, tournament *models.Tournament) (uuid.UUID, error) {
	const query = `
		INSERT INTO Tournaments(name, deposit) VALUES ($1, $2)
			RETURNING id;
	`
	var id uuid.UUID

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return id, kerror.Newf(kerror.SQLPrepareStatementError, "prepare statement %v: %w", query, err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRowContext(ctx, tournament.Name, tournament.Deposit).Scan(&id); err != nil {
		return id, kerror.Newf(kerror.SQLConstraintError, "insert tournament: %w", err)
	}

	return id, nil
}

func (tr *TournamentRepository) SelectByID(ctx context.Context, store tx.DBTX, id uuid.UUID) (*models.Tournament, error) {
	const query = `
		SELECT * FROM Tournaments WHERE id = $1
	`
	tournament := &models.Tournament{}

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt %v: %v", query, err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRowContext(ctx, id).Scan(&tournament.ID, &tournament.Name, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, kerror.Newf(kerror.TournamentDoesntExists, "tournament with id(%v) isn't exists: %v", id, err)
		}

		return nil, kerror.Newf(kerror.SQLScanError, "scan query: %v", err)
	}

	users, err := tr.selectUserIDsOfTournament(ctx, store, id)
	if err != nil {
		return nil, kerror.Errorf(err, "get users of tournament")
	}
	tournament.Users = users

	return tournament, nil

}

func (tr *TournamentRepository) selectUserIDsOfTournament(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) ([]models.User, error) {
	const query = `
		SELECT Users.id, Users.name, Users.balance
		FROM UsersOfTournaments INNER JOIN Users ON Users.id = UsersOfTournaments.id
		WHERE tournamentID = $1;
	`
	users := []models.User{}

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare query: %v", err)
	}
	defer debugutil.Close(stmt)

	rows, err := stmt.QueryContext(ctx, tournamentID)
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

func (tr *TournamentRepository) SelectRandomUserOfTournament(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) (*models.User, error) {
	const query = `
		WITH random_id AS (
			SELECT user FROM UsersOfTournaments WHERE tournamentID = $1
				OFFSET random() * COUNT(*) LIMIT 1
		)
		SELECT * FROM Users WHERE id = (SELECT id FROM random_id);
	`
	var user models.User

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return nil, kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if err := stmt.QueryRowContext(ctx, tournamentID).Scan(&user.ID, &user.Name, &user.Balance); err != nil {
		return nil, kerror.Newf(kerror.SQLScanError, "scan user from db: %v", err)
	}

	return &user, nil
}

func (tr *TournamentRepository) InsertUserToTournament(ctx context.Context, store tx.DBTX, tournamentID, userID uuid.UUID) error {
	const query = `
		INSERT INTO UsersOfTournaments(tournamentID, userID) VALUES ($1, $2); 
	`

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare query: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.ExecContext(ctx, tournamentID, userID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec stmt: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) AddToPrize(ctx context.Context, store tx.DBTX, ID uuid.UUID, d float64) error {
	const query = `
		UPDATE Tournaments
			SET prize = prize + $1
			WHERE id = $2
	`

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.ExecContext(ctx, d, ID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec query: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) RefundDepositToUsers(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID) error {
	const query = `
		WITH depositOfTournament AS (
			SELECT deposit FROM Tournaments WHERE id = $1
		)
		UPDATE Users
			SET balance = balance + (SELECT deposit FROM depositOfTournament)
			WHERE id IN (SELECT id FROM UsersOfTournaments WHERE tournamentID = $1);
	`

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.ExecContext(ctx, tournamentID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec query: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) SetWinner(ctx context.Context, store tx.DBTX, tournamentID, winnerID uuid.UUID) error {
	const query = `
		UPDATE Tournaments SET winner = $1 WHERE id = $2;
	`

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.ExecContext(ctx, winnerID, tournamentID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec update query: %v", err)
	}

	return nil
}

func (tr *TournamentRepository) UpdateStatus(ctx context.Context, store tx.DBTX, tournamentID uuid.UUID, newStatus models.TournamentStatus) error {
	const query = `
		UPDATE Tournaments SET status = $1 WHERE id = $2;
	`

	stmt, err := store.PrepareContext(ctx, query)
	if err != nil {
		return kerror.Newf(kerror.SQLPrepareStatementError, "prepare stmt: %v", err)
	}
	defer debugutil.Close(stmt)

	if _, err := stmt.ExecContext(ctx, newStatus, tournamentID); err != nil {
		return kerror.Newf(kerror.SQLExecutionError, "exec update query: %v", err)
	}

	return nil

}
