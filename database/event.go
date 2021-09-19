package database

import "virunus.com/cornhole/config"

type Event struct {
	EventId int
	Name    string
	Style   int
	Date	string
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

func (e *Event) GetGames(config config.Config) ([]*Game, error) {
	var games []*Game

	rows, err := config.Database.Query(`SELECT * FROM game WHERE eventId = ? order by gameId`, e.EventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		g := new(Game)
		if err := rows.Scan(&g.GameId, &g.EventId, &g.TeamOne, &g.TeamTwo, &g.State, &g.Winner, &g.PrevGame, &g.WinGame, &g.LoseGame); err != nil {
			return nil, err
		}

		games = append(games, g)
	}

	return games, nil
}
