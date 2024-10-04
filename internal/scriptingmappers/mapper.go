package scriptingmappers

import (
	"github.com/Ashutoshbind15/gogameengine/internal/types"
	lua "github.com/yuin/gopher-lua"
)

func MapGameStateToLuaTable(L *lua.LState, gameState *types.GameState) *lua.LTable {
	tbl := L.NewTable()

	// Enemies
	enemiesTbl := L.NewTable()
	for _, enemy := range gameState.Enemies {
		enemyTbl := L.NewTable()
		L.SetField(enemyTbl, "health", lua.LNumber(enemy.Health))
		L.SetField(enemyTbl, "px", lua.LNumber(enemy.Px))
		L.SetField(enemyTbl, "py", lua.LNumber(enemy.Py))
		enemiesTbl.Append(enemyTbl)
	}
	L.SetField(tbl, "enemies", enemiesTbl)

	// Time and timeleft
	L.SetField(tbl, "time", lua.LNumber(gameState.Time))
	L.SetField(tbl, "timeleft", lua.LNumber(gameState.TimeLeft))

	// Arena
	arenaTbl := L.NewTable()
	L.SetField(arenaTbl, "x", lua.LNumber(gameState.Arena.X))
	L.SetField(arenaTbl, "y", lua.LNumber(gameState.Arena.Y))
	L.SetField(tbl, "arena", arenaTbl)

	// Players
	playersTbl := L.NewTable()
	for _, player := range gameState.Players {
		playerTbl := L.NewTable()
		L.SetField(playerTbl, "name", lua.LString(player.Name))
		L.SetField(playerTbl, "px", lua.LNumber(player.Px))
		L.SetField(playerTbl, "py", lua.LNumber(player.Py))
		playersTbl.Append(playerTbl)
	}
	L.SetField(tbl, "players", playersTbl)

	return tbl
}

// Helper function to convert Lua table back to GameState
func MapLuaTableToGameState(L *lua.LState, tbl *lua.LTable) types.GameState {
	gameState := types.GameState{}

	// Enemies
	enemiesTbl := L.GetField(tbl, "enemies").(*lua.LTable)
	enemies := []types.Enemy{}
	enemiesTbl.ForEach(func(_, value lua.LValue) {
		enemyTbl := value.(*lua.LTable)
		enemy := types.Enemy{
			Health: int(L.GetField(enemyTbl, "health").(lua.LNumber)),
			Px:     int(L.GetField(enemyTbl, "px").(lua.LNumber)),
			Py:     int(L.GetField(enemyTbl, "py").(lua.LNumber)),
		}
		enemies = append(enemies, enemy)
	})
	gameState.Enemies = enemies

	// Time and timeleft
	gameState.Time = int(L.GetField(tbl, "time").(lua.LNumber))
	gameState.TimeLeft = int(L.GetField(tbl, "timeleft").(lua.LNumber))

	// Arena
	arenaTbl := L.GetField(tbl, "arena").(*lua.LTable)
	gameState.Arena = types.Arena{
		X: int(L.GetField(arenaTbl, "x").(lua.LNumber)),
		Y: int(L.GetField(arenaTbl, "y").(lua.LNumber)),
	}

	// Players
	playersTbl := L.GetField(tbl, "players").(*lua.LTable)
	players := []types.Player{}
	playersTbl.ForEach(func(_, value lua.LValue) {
		playerTbl := value.(*lua.LTable)
		player := types.Player{
			Name: L.GetField(playerTbl, "name").(lua.LString).String(),
			Px:   int(L.GetField(playerTbl, "px").(lua.LNumber)),
			Py:   int(L.GetField(playerTbl, "py").(lua.LNumber)),
		}
		players = append(players, player)
	})
	gameState.Players = players

	return gameState
}
