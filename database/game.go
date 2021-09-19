package database

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
