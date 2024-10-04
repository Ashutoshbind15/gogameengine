package types

type Arena struct {
	X, Y int
}

type Enemy struct {
	Health int
	Px, Py int
}

// now the dynamic funcs are described as a string of a gamescript
// eg attack -> attackfn
type Player struct {
	Name     string
	Px, Py   int
	Dynamics map[string]string
}

type GameState struct {
	Enemies  []Enemy
	Time     int
	TimeLeft int
	Arena    Arena
	Players  []Player
}
