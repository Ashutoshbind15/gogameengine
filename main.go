package main

import "fmt"

type Arena struct {
	x,y int
}

type Enemy struct {
	health int
	px, py int
}

type GameState struct {
	enemies []Enemy
	time, timeleft int
	arena Arena
	players []Player
}

type Player struct{
	name string
	px,py int
	attack func(gs GameState) GameState
}

type ResourceManager struct {
	players []Player
}


// something that the end user will be able to edit via a code interface, will be named attackinit or smth
func playerSingleLevelAttackInit (gs GameState) GameState {
	playersInGame:= gs.players

	if len(playersInGame) < 1 {
		return GameState{}
	}

	// todo: make it multiplayer
	player:= playersInGame[0]

	py := player.py
	// todo: ensure upper bounds on nextlevel
	nextlvl := py + 1

	towersInGame := gs.enemies

	// todo: use pointers instead

	resEnemies:= []Enemy{}

	for _, tower := range towersInGame {
		ty := tower.py
		if ty == nextlvl {
			chealth := tower.health
			newEnemy := Enemy {
				health: chealth-1,
				px: tower.px,
				py: tower.py,
			}
			resEnemies = append(resEnemies, newEnemy)
		}
	}


	return GameState{
		enemies: resEnemies,
		time: gs.time,
		timeleft: gs.timeleft,
		arena: gs.arena,
		players: gs.players,
	}

}

type StateChangeHook struct {
	hook func(gs GameState) GameState
	title string
}

type PlayerPlugin struct {
	stateChangeHooks []StateChangeHook
}

func (plugin *PlayerPlugin) addStateChangeHook(fn func(gs GameState) GameState, name string) {
	stateHook := StateChangeHook{
		hook: fn,
		title: name,
	}
	plugin.stateChangeHooks = append(plugin.stateChangeHooks, stateHook)
}

func test() {
	// whenever a new player is created on the server, add a pointer to it's corresponding plugin
	// load the player, load the plugin and have a mapping for the player funcs to the plugin funcs
	// i.e. player.attack() => player.plugin.attack()
}

func main() {
	fmt.Println("Game init")
}
