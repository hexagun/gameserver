package main

import (
	"fmt"
	"strconv"

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
	header := action.GetHeader()

	var im common.OutgoingMessage

	switch header.Type {
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
			Type:     "join",
			PlayerID: strconv.Itoa(header.PlayerId),
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
		panic("Unknown action to update the players with.")
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

func handleConnection(state *GameState, action common.Action) {

	header := action.GetHeader()

	if state.players[0].name == "" {
		state.players[0].name = fmt.Sprintf("%d", header.PlayerId)
		// 	state.PlayerXReady = true
		// 	state.Turn = "X"
		fmt.Printf("Player 1 (%s) has connected.\n", state.players[0].name)
	} else if state.players[1].name == "" {
		state.players[1].name = fmt.Sprintf("%d", header.PlayerId)
		// 	state.PlayerO = player
		// 	state.PlayerOReady = true
		fmt.Printf("Player 2 (%s) has connected.\n", state.players[1].name)
	} else {
		fmt.Printf("Player %s cannot join, both slots are filled.\n", header.PlayerId)
	}
}

func rootReducer(state *GameState, action common.Action) {
	switch action.GetHeader().Type {
	// case "PlayToken":
	// var playTokenAction PlayTokenAction = action.(PlayTokenAction)
	// var payload PlayTokenPayload = playTokenAction.GetPayload().(PlayTokenPayload)
	// playToken(state, payload)
	case common.Join:
		handleConnection(state, action)
	// case "DisconnectPlayer":
	// var disconnectPlayerAction DisconnectPlayerAction = action.(DisconnectPlayerAction)
	// var payload DisconnectPlayerPayload = disconnectPlayerAction.GetPayload().(DisconnectPlayerPayload)
	// handleDisconnection(state, payload)
	default:
	}
}

func (g *Game) GameLoop() {
	for {
		select {
		case action := <-g.actionChannel:

			rootReducer(&g.state, action)

			g.UpdatePlayers(action)

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
}
