package database

import (
	"log"
	"virunus.com/cornhole/config"
)

func InitializeDatabse(config *config.Config) {
	db := config.Database

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(`
create table if not exists event(
    eventId INTEGER PRIMARY KEY ,
    name TEXT NOT NULL,
    style INTEGER NOT NULL,
    date TEXT NOT NULL 
);
CREATE TABLE IF NOT EXISTS team(
    teamId INTEGER PRIMARY KEY ,
    eventId INTEGER NOT NULL,
    playerOne TEXT,
    playerTwo TEXT,
    FOREIGN KEY (eventId)
    	REFERENCES event(eventId)
    	ON DELETE CASCADE 
);
CREATE TABLE IF NOT EXISTS game(
    gameId INTEGER PRIMARY KEY ,
    eventId INTEGER NOT NULL ,
    teamOne INTEGER ,
    teamTwo INTEGER ,
    state INTEGER NOT NULL DEFAULT 0,
    winner INTEGER,
    prevGame INTEGER,
    winGame INTEGER,
    loseGame INTEGER,
    FOREIGN KEY (eventId) REFERENCES event(eventId) ON DELETE CASCADE ,
    FOREIGN KEY (teamOne) REFERENCES team(teamId) ON DELETE SET NULL ,
    FOREIGN KEY (teamTwo) REFERENCES team(teamId) ON DELETE SET NULL ,
    FOREIGN KEY (winner) REFERENCES team(teamId) ON DELETE SET NULL
);
`)
	if err != nil {
		log.Print(err.Error())

		err := tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
}

func GetEvents(config *config.Config) ([]Event, error) {
	var events []Event

	rows, err := config.Database.Query(`SELECT * FROM event ORDER BY date DESC, name LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.EventId, &e.Name, &e.Style, &e.Date); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (e *Event) GetTeams(config *config.Config) ([]Team, error) {
	var teams []Team

	rows, err := config.Database.Query(`SELECT * FROM team WHERE eventId = ? ORDER BY teamId;`, e.EventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.TeamId, &t.EventId, &t.PlayerOne, &t.PlayerTwo); err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}

	return teams, nil
}
