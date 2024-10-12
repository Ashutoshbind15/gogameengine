package realtime

import (
	"time"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
)

// todo: export the types from the types package

type GameInstanceResources struct {
	Players []types.Player
	Time time.Time
	TimeStarted time.Time
}

// todo: separate out a gameinstance for the dynamic/realtime parts of the game to handle in mem, from the metadata

type Game struct {
	Id string
	InstanceId string
	Gamestate *types.GameState
	Clients []*GameClient
	Moves []types.GameAction // just in case if we want to reconstruct the state from the moves
	GameInfo *GameMeta
	InstanceResources *GameInstanceResources
	Aggregator chan []byte
	TurnBitmap string
}

func (gm *Game) Runner() {
	for bcastmsg := range gm.Aggregator {
		for _, client := range gm.Clients {
			client.Send <- bcastmsg
		}
	}
}