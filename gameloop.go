package main

import (
	"fmt"

	"github.com/hexagun/common"
	"github.com/hexagun/gameserver/websocket"
)

type Player struct {
	id    int
	name  string
	token string
}

type GameState struct {
	board             [3][3]string
	players           [2]Player
	activePlayerIndex int
	winner            string
	draw              bool
}

//type GameRenderer func(game *Game)

type Game struct {
	state         GameState
	actionChannel chan common.Action
	pool          *websocket.Pool // Seems wrong
}

// func (c *Client) DecodeMessage(msg *common.IncomingMessage) {
// 	c.Decoder(msg)
// }

func (g *Game) Dispatch(action common.Action) {
	g.actionChannel <- action
}

func (g *Game) UpdatePlayers(action common.Action) {
	actionType := action.GetHeader().Type

	var im common.OutgoingMessage

	switch actionType {
	case common.GameStateUpdate:
		im = common.OutgoingMessage{
			Type:    "gamestateupdate",
			Payload: action.GetPayload(),
		}
	case common.GameOver:
		im = common.OutgoingMessage{
			Type:    "gameover",
			Payload: action.GetPayload(),
		}
	case common.Join:
		im = common.OutgoingMessage{
			Type: "join",
		}
	case common.Leave:
		im = common.OutgoingMessage{
			Type: "leave",
		}
	case common.Start:
		im = common.OutgoingMessage{
			Type: "start",
		}
	case common.Error:
		im = common.OutgoingMessage{
			Type:    "error",
			Payload: action.GetPayload(),
		}
	default:
		panic("Unknown action to update the players with")
	}

	g.pool.Broadcast <- im

}

func NewGame(pool *websocket.Pool) *Game {
	return &Game{
		state:         GameState{},
		actionChannel: make(chan common.Action),
		pool:          pool,
	}
}

func (g *Game) GameLoop() {
	select {
	case action := <-g.actionChannel:
		g.UpdatePlayers(action)
		fmt.Println(action)
		// // // Handle the action
		// rootReducer(&state, action)
		// printBoard(state.board)

		// if state.winner != "" {
		// 	fmt.Printf("Player %s wins!\n", state.winner)
		// 	//return
		// } else if state.draw {
		// 	fmt.Println("It's a draw!")
		// 	//return
		// }

		// // Proceed to the next turn
		// if state.activePlayerIndex == 0 {
		// 	state.activePlayerIndex = 1
		// } else {
		// 	state.activePlayerIndex = 0
		// }
	}
}
