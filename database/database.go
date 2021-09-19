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

// GetEvents returns a slice containing a list of all events from the database
func GetEvents(config *config.Config) ([]*Event, error) {
	var events []*Event

	rows, err := config.Database.Query(`SELECT * FROM event ORDER BY date DESC, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := new(Event)
		if err := rows.Scan(&e.EventId, &e.Name, &e.Style, &e.Date); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

// Save saves the current state of the event to the database, creating a new entry if necessary
func (e *Event) Save(config *config.Config) (*Event, error) {
	if e.EventId == 0 {
		_, err := config.Database.Exec(
			"INSERT INTO event(name, style, date) VALUES (?, ?, ?);",
			e.Name,
			e.Style,
			e.Date,
		)
		if err != nil {
			return nil, err
		}

		row := config.Database.QueryRow("SELECT last_insert_rowid() FROM event;")

		err = row.Scan(&e.EventId)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := config.Database.Exec(
			"UPDATE event SET name = ?, style = ?, date = ? WHERE eventId = ?",
			e.Name,
			e.Style,
			e.Date,
			e.EventId,
		)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

// Delete removes the event and associated teams and games from the database
func (e *Event) Delete(config *config.Config) error {
	// todo: explicitly delete associated objects instead of relying on constraints
	_, err := config.Database.Exec("DELETE FROM event WHERE eventId = ?", e.EventId)

	return err
}

// GetTeams retrieves all teams associated with the event as a slice
func (e *Event) GetTeams(config *config.Config) ([]*Team, error) {
	var teams []*Team

	rows, err := config.Database.Query(`SELECT * FROM team WHERE eventId = ? ORDER BY teamId;`, e.EventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := new(Team)
		if err := rows.Scan(&t.TeamId, &t.EventId, &t.PlayerOne, &t.PlayerTwo); err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}

	return teams, nil
}

// Save saves the current state of the event to the database, creating a new entry if necessary
func (t *Team) Save(config *config.Config) (*Team, error) {
	if t.TeamId == 0 {
		_, err := config.Database.Exec(
			"INSERT INTO team(eventId, playerOne, playerTwo) VALUES (?, ?, ?);",
			t.EventId,
			t.PlayerOne,
			t.PlayerTwo,
		)
		if err != nil {
			return nil, err
		}

		row := config.Database.QueryRow("SELECT last_insert_rowid() FROM team;")

		err = row.Scan(&t.TeamId)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := config.Database.Exec(
			"UPDATE team SET playerOne = ?, playerTwo = ? WHERE teamId = ?;",
			t.PlayerOne,
			t.PlayerTwo,
			t.TeamId,
		)
		if err != nil {
			return nil, err
		}
	}


	return t, nil
}

// Delete removes the team from the database
func (t *Team) Delete(config *config.Config) error {
	_, err := config.Database.Exec("DELETE FROM team WHERE teamId = ?", t.TeamId)

	return err
}
