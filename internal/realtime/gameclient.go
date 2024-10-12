package realtime

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type GameClient struct {
	Connection *websocket.Conn
	CurrentGame *Game
	CurrentGameId string
	Cidx int
	UserId uuid.UUID
	Username string
	IsWaiting bool
	IsConnected bool
	Send chan []byte
	SendErr chan []byte
}

type ClientRequest struct {
	ActionOpCode string
	Data string
}

type JoinRelay struct {
	GameToJoin *Game
	ClientToJoin *GameClient
}

func isAMove(opcode string) bool {
	return opcode == "action" 
}

func (gc *GameClient) Reader(handlercloser chan bool) {

	// Listens to the ws msgs and relays them to the manager to process, and send back to the writers
	defer func() {
		Manager.RemoveClient <- gc.Connection
		gc.Connection.Close()
		// todo: check if it's okay
		gc.Connection = nil
		handlercloser <- true
	}()

	buff := make([]byte, 1024)

	for {		
		n, err := gc.Connection.Read(buff)


		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed for: ", gc.UserId)
				return;
			}

			fmt.Println("ERR IN READING DATA FOR CLIENT: ", gc.UserId, err)
			continue;
			// break;
		}

		var rq ClientRequest;
		res := buff[:n]
		json.Unmarshal(res, &rq)

		fmt.Println(rq)

		if rq.ActionOpCode == "join" {

			res := strings.Split(rq.Data, ":")
			roomId := res[1]

			// todo:doubt read can be done by multiple go routines too even if concurrent?	
			game, err := Manager.FindGameByInstanceId(roomId)

			if err != nil {
				fmt.Println("Cannot find game with the id: ", roomId, err)
				continue
			}

			jrr := JoinRelay{
				GameToJoin: game,
				ClientToJoin: gc,
			}

			Manager.JoinGame <- jrr

		} else if rq.ActionOpCode == "create" {
			
			res := strings.Split(rq.Data, ":")
			gameId, roomId := res[0], res[1]
			
			ng := NewGame{
				GameId: gameId,
				GameInstanceId: roomId,
			}

			ngb, err := json.Marshal(ng)

			if err != nil {
				fmt.Println("Error in creating the byte info for adding the new game: ", err)
				continue;
			}

			Manager.Creategame <- ngb

		} else if rq.ActionOpCode == "start" {
			gameInstanceId := rq.Data
			Manager.Startgame <- gameInstanceId
		} else if rq.ActionOpCode == "gamemove" {
			res := strings.Split(rq.Data, ":")
			gameInstanceId := res[0]
			game, err := Manager.FindGameByInstanceId(gameInstanceId)

			if err != nil {
				fmt.Println("Cannot find the game with the gameId: ", gameInstanceId, err)
				continue;
			}

			game.Aggregator <- []byte(res[1])
		}
		
	}
}

func (gc *GameClient) Sender(handlercloser chan bool) {

	defer func() {
		// todo: see if the closer from two ends (w/r) dont led to blocking as we're listening on the handler only once, instead, use a non blocking send
		handlercloser <- true
	}()

	for  msg := range gc.Send {
		// todo: check whether the close has been called on the ws and the channel hasn't been closed 

		if gc.Connection == nil {
			break;
		}

		gc.Connection.Write(msg)
	}
}