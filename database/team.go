package database

import "virunus.com/cornhole/config"

type Team struct {
	TeamId    int
	EventId   int
	PlayerOne string
	PlayerTwo string
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
