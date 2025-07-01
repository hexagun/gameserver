package main

import (
	"fmt"
	"math/rand"
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
	state             State
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

func (g *Game) Encode(action common.Action) common.OutgoingMessage {
	var im common.OutgoingMessage
	header := action.GetHeader()
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
			GameID:   strconv.Itoa(header.GameId),
		}
	case common.Leave:
		im = common.OutgoingMessage{
			Type: "leave",
		}
	case common.Start:
		im = common.OutgoingMessage{
			Type:     "start",
			PlayerID: strconv.Itoa(header.PlayerId),
			GameID:   strconv.Itoa(header.GameId),
			Payload:  action.GetPayload(),
		}
	case common.Error:
		im = common.OutgoingMessage{
			Type:    "error",
			Payload: action.GetPayload(),
		}
	default:
		panic("Unknown action to update the players with.")
	}
	return im
}

func (g *Game) UpdatePlayer(id string, action common.Action) {
	im := g.Encode(action)
	g.pool.Send <- im
}

func (g *Game) UpdatePlayers(actionList []common.Action) {
	for _, action := range actionList {
		//header := action.GetHeader()
		im := g.Encode(action)
		g.pool.Broadcast <- im
	}
}

func NewGame(pool *websocket.Pool) *Game {
	return &Game{
		state: GameState{
			state: InitialState,
		},
		actionChannel: make(chan common.Action),
		pool:          pool,
	}
}

func handleConnection(state *GameState, action common.Action) bool {

	header := action.GetHeader()

	if state.players[0].name == "" {
		state.players[0].name = fmt.Sprintf("%d", header.PlayerId)
		tr := rules[state.state][0] // Join
		state.state = tr.State
		fmt.Printf("Player 1 (%s) has connected.\n", state.players[0].name)
	} else if state.players[1].name == "" {
		state.players[1].name = fmt.Sprintf("%d", header.PlayerId)
		tr := rules[state.state][0] // Join
		state.state = tr.State
		fmt.Printf("Player 2 (%s) has connected.\n", state.players[1].name)
	} else {
		fmt.Printf("Player %s cannot join, both slots are filled.\n", header.PlayerId)
		return false
	}

	return true
}

func handleDisconnection(state *GameState, action common.Action) {

	header := action.GetHeader()
	if state.players[0].name == strconv.Itoa(header.PlayerId) {
		state.players[0].name = ""
		// 	state.PlayerXReady = true
		// 	state.Turn = "X"
		tr := rules[state.state][Leave]
		state.state = tr.State
		fmt.Printf("Player 1 (%s) has disconnected.\n", header.PlayerId)
	} else if state.players[1].name == "" {
		state.players[1].name = fmt.Sprintf("%d", header.PlayerId)
		// 	state.PlayerO = player
		// 	state.PlayerOReady = true
		tr := rules[state.state][Leave]
		state.state = tr.State
		fmt.Printf("Player 2 (%s) has disconnected.\n", state.players[1].name)
	} else {
		fmt.Printf("Player %s cannot disconnect, both slots are filled.\n", header.PlayerId)
	}

}

func startGame(state *GameState, action common.Action) {

	tr := rules[state.state][0] // Start
	state.state = tr.State

	// header := action.GetHeader()

	// return common.NewGameStateUpdateAction(header.GameId, 0,
	// 	common.GameStateUpdatePayload{
	// 		Board:    state.board,
	// 		NextTurn: state.players[state.activePlayerIndex].name,
	// 		Winner:   "", // will never be called
	// 	})
}

func checkWin(board [3][3]string, player string) bool {
	for i := 0; i < 3; i++ {
		// Check rows and columns
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}
	// Check diagonals
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}
	return false
}

// Check if the game ended in a draw
func checkDraw(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func playerMove(state *GameState, action common.Action) common.Action {
	header := action.GetHeader()
	payload := action.GetPayload().(common.PlayerMovePayload)

	if state.board[payload.Col][payload.Row] != "" {
		fmt.Println("Invalid move! Cell already taken.")
		return nil // error action
	}

	var token string

	if state.players[0].name == strconv.Itoa(header.PlayerId) {
		// Place the player's mark on the board
		token = state.players[0].token
	} else {
		token = state.players[1].token
	}

	state.board[payload.Col][payload.Row] = token

	// Check for a win or draw after the move
	if checkWin(state.board, token) {
		//state.winner = payload.Player
		gaAction := common.NewGameOverAction(header.GameId,
			common.GameOverPayload{
				Winner: strconv.Itoa(header.PlayerId),
				Board:  state.board,
			})
		tr := rules[state.state][2] //GameEnded
		state.state = tr.State
		return gaAction

	} else if checkDraw(state.board) {
		gaAction := common.NewGameOverAction(header.GameId,
			common.GameOverPayload{
				Board: state.board,
			})
		tr := rules[state.state][2] //GameEnded
		state.state = tr.State
		return gaAction
	}

	// Proceed to the next turn
	if state.activePlayerIndex == 0 {
		state.activePlayerIndex = 1
	} else {
		state.activePlayerIndex = 0
	}

	return common.NewGameStateUpdateAction(header.GameId, 0,
		common.GameStateUpdatePayload{
			Board:    state.board,
			NextTurn: state.players[state.activePlayerIndex].name,
			Winner:   "", // will never be called
		})

}

func initialStateReducer(state *GameState, action common.Action) common.Action {
	switch action.GetHeader().Type {
	case common.Join:
		handleConnection(state, action)
	default:
	}
	return nil
}

func coinToss() int {
	return rand.Int() % 2
}

func waitingForStateReducer(state *GameState, action common.Action) common.Action {
	switch action.GetHeader().Type {
	case common.Join:
		if handleConnection(state, action) {
			// type StartPayload struct {
			// 	YourToken  string //`json:"yourToken"`
			// 	OpponentID string //`json:"opponentId"`
			// 	FirstTurn  string //`json:"firstTurn"` // could also be a player ID
			// }

			//coin toss for token
			startingPlayerToken := "o"
			opponentPlayerToken := "x"
			if coinToss() == 0 {
				startingPlayerToken = "x"
				opponentPlayerToken = "o"
			}

			startingPlayerIndex := coinToss()
			state.activePlayerIndex = startingPlayerIndex
			opponentPlayerIndex := 0
			if startingPlayerIndex == 0 {
				opponentPlayerIndex = 1
			}

			state.players[startingPlayerIndex].token = startingPlayerToken
			state.players[opponentPlayerIndex].token = opponentPlayerToken

			playerId, _ := strconv.Atoi(state.players[startingPlayerIndex].name)

			// Starting player action
			return common.NewStartAction(action.GetHeader().GameId,
				playerId,
				common.StartPayload{
					YourToken:  startingPlayerToken,
					OpponentID: state.players[opponentPlayerIndex].name,
					FirstTurn:  state.players[startingPlayerIndex].name,
				})
		}

	case common.Leave:
		handleDisconnection(state, action)
	default:
	}

	return nil
}
func readyReducer(state *GameState, action common.Action) common.Action {
	switch action.GetHeader().Type {
	case common.Start:
		startGame(state, action)
		return action
	case common.Leave:
		handleDisconnection(state, action)
	default:
	}
	return nil
}
func inGameReducer(state *GameState, action common.Action) common.Action {

	switch action.GetHeader().Type {
	case common.Move:
		return playerMove(state, action)
	case common.Leave:
		handleDisconnection(state, action)
	default:
	}
	return nil
}
func rejoinReducer(state *GameState, action common.Action) {
	switch action.GetHeader().Type {
	case common.Join:
		handleConnection(state, action)
	case common.Leave:
		handleDisconnection(state, action)
	default:
	}
}
func gameEndedReducer(state *GameState, action common.Action) {
	switch action.GetHeader().Type {
	case common.Leave:
		handleDisconnection(state, action)
	default:
	}
}
func (g *Game) GameLoop() {

	for {
		actionList := make([]common.Action, 0)
		select {
		case action := <-g.actionChannel:
			broadcast := true
			fmt.Printf("Processing action (%s)\n", action)
			switch g.state.state {
			case InitialState:
				initialStateReducer(&g.state, action)
				actionList = append(actionList, action)
			case WaitingForOpponent:
				broadcast = false
				startAction := waitingForStateReducer(&g.state, action)
				if startAction != nil {
					go func() {
						g.Dispatch(startAction)
					}()
					actionList = append(actionList, action) // join action
				}
			case ReadyToStart:
				startAction := readyReducer(&g.state, action)
				if startAction != nil {

					actionList = append(actionList, startAction)

					// Find opponent ID
					startingPlayerIndex := g.state.activePlayerIndex
					opponentPlayerIndex := 0
					if startingPlayerIndex == 0 {
						opponentPlayerIndex = 1
					}

					opponentId, _ := strconv.Atoi(g.state.players[opponentPlayerIndex].name)

					// Starting player action
					opponentStartAction := common.NewStartAction(action.GetHeader().GameId,
						opponentId,
						common.StartPayload{
							YourToken:  g.state.players[opponentPlayerIndex].token,
							OpponentID: g.state.players[startingPlayerIndex].name,
							FirstTurn:  g.state.players[startingPlayerIndex].name,
						})

					actionList = append(actionList, opponentStartAction)

					broadcast = false // Send individual messages
				}

			case InGame:
				gameAction := inGameReducer(&g.state, action)
				if gameAction != nil {
					actionList = append(actionList, gameAction)
					go func() {
						g.Dispatch(gameAction)
					}()
				}
			case Rejoin:
				rejoinReducer(&g.state, action)
			case GameEnded:
				gameEndedReducer(&g.state, action)
			default:
				//case State.GameOver:
				//case State.DoubleDisconnect:
				//case State.FinalState:
				// quit server
			}

			// multiple Response action, 2 max problaby (each player)
			// for example, startAction

			// if game started, collect play action and deduce if win
			// Draw, kind of
			if broadcast {
				g.UpdatePlayers(actionList)
			} else {
				for _, action := range actionList {
					g.UpdatePlayer(strconv.Itoa(action.GetHeader().PlayerId), action)
				}
			}

			// // // Handle the action
			// rootReducer(&state, action)o
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
