package realtime

import (
	"encoding/json"
	"fmt"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
)

type GameMeta struct {
	Id string
	Title string
	Description string
	Url string
	Players []types.Player
	Resources []GameInstanceResources
}

type GameManager struct {
	Games     []Game
	Closegame chan bool
	Startgame chan []byte
	LoadedGamesMetadata map[string]*GameMeta // id -> metadata for the current games
}

type NewGame struct {
	GameId       string
	GameInstanceId string
}

func (gm *GameManager) loadGamesMetadata(id string) *GameMeta {
	// todo: make a db call if the game meta isn't there in mem
	return gm.LoadedGamesMetadata[id]
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
				TurnBitmap: "",
				GameInfo: gm.loadGamesMetadata(gameInfo.GameId),
			}
			gm.Games = append(gm.Games, tempGame)
			
		case closeFlag :=  <- gm.Closegame:
			// close the game, deallocate the resources
			fmt.Println("CLOSING THE GAME: ", closeFlag)
		}
	}
}