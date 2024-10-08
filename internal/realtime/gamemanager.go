package realtime

import (
	"encoding/json"
	"fmt"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
)

type GameManager struct {
	Games     []Game
	Closegame chan bool
	Startgame chan []byte
}

type NewGame struct {
	GameId       string
	GameInstanceId string
}

func (gm *GameManager) Manage() {
	for {
		select {
		case newGameInfo := <-gm.Startgame:

			var gameInfo NewGame

			json.Unmarshal(newGameInfo, &gameInfo)

			// todo: find the game details and populate the initialization here using the gameid

			tempGame := Game{
				InstanceId: gameInfo.GameInstanceId,
				Id: gameInfo.GameId,
				Clients: []GameClient{},
				Gamestate: &types.GameState{},
				Moves: []types.GameAction{},
				BroadCast: make(chan []byte),
				ClientAction: make(chan []byte),
			}
			gm.Games = append(gm.Games, tempGame)
			
		case closeFlag :=  <- gm.Closegame:
			// close the game, deallocate the resources
			fmt.Println("CLOSING THE GAME: ", closeFlag)
		}
	}
}