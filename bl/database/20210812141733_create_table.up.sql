 CREATE TABLE IF NOT EXISTS Users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name varchar(200) NOT NULL UNIQUE,
	balance numeric(10, 2) NOT NULL DEFAULT 0 CHECK(balance >= 0.0)
);

CREATE TYPE TournamentStatus AS ENUM ('Active', 'Cancel', 'Finish');

CREATE TABLE IF NOT EXISTS Tournaments (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name varchar(200) NOT NULL,
	deposit numeric(10, 2) NULL DEFAULT 0 CHECK(deposit >= 0.0),
	prize  numeric(12, 2) NOT NULL DEFAULT 0 CHECK(prize >= 0.0),
	winner uuid REFERENCES Users(id),
	status TournamentStatus NOT NULL DEFAULT 'Active'
);

CREATE TABLE IF NOT EXISTS UsersOfTournaments (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	tournamentID uuid REFERENCES Tournaments(id) NOT NULL,
	userID uuid REFERENCES Users(id) NOT NULL
);


