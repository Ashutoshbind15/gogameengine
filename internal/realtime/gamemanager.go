package realtime

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

var (
	Manager *GameManager
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
	Games     []*Game
	Closegame chan []byte
	Startgame chan string
	Creategame chan []byte
	JoinGame chan JoinRelay

	RemoveClient chan *websocket.Conn
	AddClient chan NewClientInfo
	LoadedGamesMetadata map[string]*GameMeta // id -> metadata for the current games
	ConnectionMeta map[*websocket.Conn]*GameClient
}

type NewGame struct {
	GameId       string
	GameInstanceId string
}

type NewClientInfo struct {
	UserId uuid.UUID
	Username string
	Conn *websocket.Conn
	CloseChan chan bool
}

func (gm *GameManager) loadGamesMetadata(id string) *GameMeta {
	// todo: make a db call if the game meta isn't there in mem
	
	time.Sleep(time.Second)
	return gm.LoadedGamesMetadata[id]
}

func (gm *GameManager) FindGameByInstanceId (id string) (*Game, error) {
	for _, gam := range gm.Games {
		if gam.InstanceId == id {
			return gam, nil
		}
	}

	return &Game{}, fmt.Errorf("cannot find a game with that instance id")
}

func (gm *GameManager) Manage() {
	for {
		select {
		case newClient := <- gm.AddClient:

			client := &GameClient{
				Connection: newClient.Conn,
				UserId: newClient.UserId,
				Send: make(chan []byte),
				SendErr: make(chan []byte),
				IsConnected: true,
				IsWaiting: true,
				Username: newClient.Username,
			}
			
			Manager.ConnectionMeta[newClient.Conn] = client

			go client.Reader(newClient.CloseChan)
			go client.Sender(newClient.CloseChan)		

		case deleteClient := <-gm.RemoveClient:
			close(gm.ConnectionMeta[deleteClient].Send)
			close(gm.ConnectionMeta[deleteClient].SendErr)
			delete(gm.ConnectionMeta, deleteClient)

		case newGameInfo := <-gm.Creategame:

			var gameInfo NewGame
			json.Unmarshal(newGameInfo, &gameInfo)

			// if gm.LoadedGamesMetadata[gameInfo.GameId] == nil {
			// 	gm.loadGamesMetadata(gameInfo.GameId)
			// }

			// todo: find the game details and populate the initialization here using the gameid

			tempGame := &Game{
				InstanceId: gameInfo.GameInstanceId,
				Id: gameInfo.GameId,
				Clients: []*GameClient{},
				Gamestate: &types.GameState{},
				Moves: []types.GameAction{},
				Aggregator: make(chan []byte),
				TurnBitmap: "",
				// GameInfo: gm.loadGamesMetadata(gameInfo.GameId),
				GameInfo: &GameMeta{},
			}
			gm.Games = append(gm.Games, tempGame)
			go tempGame.Runner()
		
		case joinGameInfo := <- gm.JoinGame:
			joinGameInfo.GameToJoin.Clients = append(joinGameInfo.GameToJoin.Clients, joinGameInfo.ClientToJoin)

		case closeFlag :=  <- gm.Closegame:
			// close the game, deallocate the resources
			fmt.Println("CLOSING THE GAME: ", closeFlag)
		case gameInstanceId := <- gm.Startgame:
			fmt.Println("Starting the game: ", gameInstanceId)
			
			game, err := gm.FindGameByInstanceId(gameInstanceId)

			if err != nil {
				fmt.Println("Cannot find the game with the gameId: ", gameInstanceId, err)
			}

			game.Aggregator <- []byte(gameInstanceId)

		}

	}
}