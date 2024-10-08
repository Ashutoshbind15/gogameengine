package extrascripts

import (
	"fmt"

	"github.com/Ashutoshbind15/gogameengine/internal/scriptingmappers"
	"github.com/Ashutoshbind15/gogameengine/internal/types"
	lua "github.com/yuin/gopher-lua"
)

func TestScript() {
	L := lua.NewState()
	defer L.Close()

	fratk := FrontalAttackScript
	
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
			{Name: "Player1", Px: 2, Py: 3, Dynamics: map[string]string{"attack": fratk}},
		},
	}

	modifiedGameState := scriptingmappers.GameStateScriptRunner(L, fratk, "attack", gameState)

	fmt.Printf("Modified GameState: %+v\n", modifiedGameState)

}