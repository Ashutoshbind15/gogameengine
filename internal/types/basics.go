package types

import "time"

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

type Session struct {
	Id      string
	UserId int
	ValidTo time.Time
}

type User struct {
	Id int
	Username string
	Password string
}

type GameAction struct {
	ActionOpCode string
	ActionData   string // this could be a json string representing the coordinates or some extra info related to the action
}