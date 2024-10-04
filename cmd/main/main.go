package main

import (
	"fmt"

	"github.com/Ashutoshbind15/gogameengine/internal/extrascripts"
	"github.com/Ashutoshbind15/gogameengine/internal/scriptingmappers"
	"github.com/Ashutoshbind15/gogameengine/internal/types"
	lua "github.com/yuin/gopher-lua"
)


func test() {
	// whenever a new player is created on the server, add a pointer to it's corresponding plugin
	// load the player, load the plugin and have a mapping for the player funcs to the plugin funcs
	// i.e. player.attack() => player.plugin.attack()

	L := lua.NewState()
	defer L.Close()

	fratk := extrascripts.FrontalAttackScript
	err := L.DoString(fratk)
	
	if err != nil {
		panic(err)
	}

	// Sample GameState
	gameState := &types.GameState{
		Enemies: []types.Enemy{
			{Health: 100, Px: 2, Py: 3},
			{Health: 80, Px: 2, Py: 4},
		},
		Time:     120,
		TimeLeft: 60,
		Arena:    types.Arena{X: 100, Y: 200},
		Players: []types.Player{
			{Name: "Player1", Px: 2, Py: 3},
		},
	}

	luaGameState := scriptingmappers.MapGameStateToLuaTable(L, gameState)

	fn := L.GetGlobal("playerSingleLevelAttackInit")
	if err := L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, luaGameState); err != nil {
		panic(err)
	}

	modifiedLuaGameState := L.Get(-1).(*lua.LTable)
	L.Pop(1)

	modifiedGameState := scriptingmappers.MapLuaTableToGameState(L, modifiedLuaGameState)

	fmt.Printf("Modified GameState: %+v\n", modifiedGameState)

}

func main() {
	fmt.Println("Game init")
	test()
}
