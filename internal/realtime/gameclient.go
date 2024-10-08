package realtime

import (
	"fmt"
	"io"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
	"golang.org/x/net/websocket"
)

type GameClient struct {
	Connection *websocket.Conn
	CurrentGame *Game
	User types.User
	Send chan []byte
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

		res := buff[:n]
		gc.CurrentGame.BroadCast <- res
	}
}

func (gc *GameClient) Sender() {
	for  msg := range gc.Send {
		gc.Connection.Write(msg)
	}
}