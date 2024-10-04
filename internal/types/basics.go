package types

type Arena struct {
	X, Y int
}

type Enemy struct {
	Health int
	Px, Py int
}

type GameState struct {
	Enemies  []Enemy
	Time     int
	TimeLeft int
	Arena    Arena
	Players  []Player
}

type Player struct {
	Name     string
	Px, Py   int
	Attack   func(gs GameState) GameState
	Dynamics []StateChanger
}

type ResourceManager struct {
	Players []Player
}

type StateChanger map[string]func(gs GameState) GameState
