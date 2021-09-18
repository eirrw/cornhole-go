package database

type Event struct {
	EventId int
	Name    string
	Style   int
	Date	string
}

type Game struct {
	GameId   int
	EventId  int
	TeamOne  int
	TeamTwo  int
	State    int
	Winner   int
	PrevGame int
	WinGame  int
	LoseGame int
}

type Team struct {
	TeamId    int
	EventId   int
	PlayerOne string
	PlayerTwo string
}