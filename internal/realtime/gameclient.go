package realtime

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
	"golang.org/x/net/websocket"
)

type GameClient struct {
	Connection *websocket.Conn
	CurrentGame *Game
	Cidx int
	User types.User
	Send chan []byte
	SendErr chan []byte
}

type ClientRequest struct {
	ActionOpCode string
	Data string
}

func isAMove(opcode string) bool {
	return opcode == "action" 
}

func (gc *GameClient) Reader() {
	buff := make([]byte, 1024)
	for {
		n, err := gc.Connection.Read(buff)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed for: ", gc.User.Id)
				return;
			}

			fmt.Println("ERR IN READING DATA FOR CLIENT: ", gc.User.Id)
			continue;
		}

		var rq ClientRequest;
		res := buff[:n]

		json.Unmarshal(res, &rq)

		if isAMove(rq.ActionOpCode) && gc.CurrentGame.TurnBitmap[gc.Cidx] == '0' {
			gc.SendErr <- []byte("Not allowed, wait for your turn")
			return
		}

		gc.CurrentGame.BroadCast <- res

	}
}

func (gc *GameClient) Sender() {
	for  msg := range gc.Send {
		gc.Connection.Write(msg)
	}
}