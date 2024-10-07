package main

import (
	"fmt"
	"net/http"

	"github.com/Ashutoshbind15/gogameengine/internal/data"
	"github.com/Ashutoshbind15/gogameengine/internal/extrascripts"
	"github.com/Ashutoshbind15/gogameengine/internal/scriptingmappers"
	"github.com/Ashutoshbind15/gogameengine/internal/types"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	lua "github.com/yuin/gopher-lua"
)


func test() {
	L := lua.NewState()
	defer L.Close()

	fratk := extrascripts.FrontalAttackScript
	
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

func initHandler (w http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(w, "init")
}


func main() {
	
	fmt.Println("Game init")

	// test()

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	data.InitDB()

	defer func(){
		data.DbConn.Close()
	}()

	// data.InitTables()

	r := mux.NewRouter()
	r.HandleFunc("/", initHandler)
	http.ListenAndServe(":3000", r)
}
