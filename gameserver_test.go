package gameserver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

var actionChannel = make(chan Action)

var state = GameState{}

func playToken(state *GameState, payload PlayTokenPayload) {
	if state.board[payload.X][payload.Y] != "" {
		fmt.Println("Invalid move! Cell already taken.")
		return
	}

	// Place the player's mark on the board
	state.board[payload.X][payload.Y] = payload.Player

	// Check for a win or draw after the move
	if checkWin(state.board, payload.Player) {
		state.winner = payload.Player
	} else if checkDraw(state.board) {
		state.draw = true
	}
}

func handleConnection(state *GameState, payload ConnectPlayerPayload) {
	if state.players[0].name == "" {
		state.players[0].name = payload.Player
		// 	state.PlayerXReady = true
		// 	state.Turn = "X"
		fmt.Printf("Player 1 (%s) has connected.\n", payload.Player)
	} else if state.players[1].name == "" {
		state.players[1].name = payload.Player
		// 	state.PlayerO = player
		// 	state.PlayerOReady = true
		fmt.Printf("Player 2 (%s) has connected.\n", payload.Player)
	} else {
		fmt.Printf("Player %s cannot join, both slots are filled.\n", payload.Player)
	}
}

func handleDisconnection(state *GameState, payload DisconnectPlayerPayload) {
	if payload.Player == state.players[0].name {
		fmt.Printf("Player X (%s) has disconnected.\n", payload.Player)
		// state.PlayerXReady = false
		// state.PlayerX = "" // Allow reconnection or replacement
	} else if payload.Player == state.players[1].name {
		fmt.Printf("Player O (%s) has disconnected.\n", payload.Player)
		// state.PlayerOReady = false
		// state.PlayerO = "" // Allow reconnection or replacement
	}
}

func rootReducer(state *GameState, action Action) {
	switch action.GetType() {
	case "PlayToken":
		var playTokenAction PlayTokenAction = action.(PlayTokenAction)
		var payload PlayTokenPayload = playTokenAction.GetPayload().(PlayTokenPayload)
		playToken(state, payload)
	case "ConnectPlayer":
		var connectPlayerAction ConnectPlayerAction = action.(ConnectPlayerAction)
		var payload ConnectPlayerPayload = connectPlayerAction.GetPayload().(ConnectPlayerPayload)
		handleConnection(state, payload)
	case "DisconnectPlayer":
		var disconnectPlayerAction DisconnectPlayerAction = action.(DisconnectPlayerAction)
		var payload DisconnectPlayerPayload = disconnectPlayerAction.GetPayload().(DisconnectPlayerPayload)
		handleDisconnection(state, payload)
	default:
	}
}

func gameLoop() {
	state.activePlayerIndex = 0 // X starts first
	state.players[0] = Player{id: 0, name: "", token: "x"}
	state.players[1] = Player{id: 1, name: "", token: "o"}

	for {
		select {
		case action := <-actionChannel:
			// // Handle the action
			rootReducer(&state, action)
			printBoard(state.board)

			if state.winner != "" {
				fmt.Printf("Player %s wins!\n", state.winner)
				//return
			} else if state.draw {
				fmt.Println("It's a draw!")
				//return
			}

			// Proceed to the next turn
			if state.activePlayerIndex == 0 {
				state.activePlayerIndex = 1
			} else {
				state.activePlayerIndex = 0
			}
		}
	}
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

// Dispatch an action to the game loop via the channel
func dispatch(action Action) {
	actionChannel <- action
}

// Utility function to print the board
func printBoard(board [3][3]string) {
	fmt.Println("Current board:")
	for _, row := range board {
		fmt.Println(row)
	}
	fmt.Println()
}

func TestGame(t *testing.T) {
	// Start the game loop in a separate goroutine
	go gameLoop()

	// Simulate player actions with dispatches
	time.Sleep(1 * time.Second) // Allow game loop to start

	dispatch(CreateConnectPlayerAction("X"))

	// Simulate player actions with dispatches
	time.Sleep(1 * time.Second) // Allow game loop to start

	dispatch(CreateConnectPlayerAction("O"))

	// Simulate player actions with dispatches
	time.Sleep(1 * time.Second) // Allow game loop to start

	dispatch(CreatePlayTokenAction("X", 0, 0))
	// Simulate player actions with dispatches
	time.Sleep(1 * time.Second) // Small delay for demonstration

	dispatch(CreatePlayTokenAction("O", 1, 1)) // O plays
	time.Sleep(1 * time.Second)

	dispatch(CreatePlayTokenAction("X", 0, 1)) // X plays
	time.Sleep(1 * time.Second)

	dispatch(CreatePlayTokenAction("O", 1, 2)) // O plays
	time.Sleep(1 * time.Second)

	dispatch(CreatePlayTokenAction("X", 0, 2)) // X wins
	time.Sleep(1 * time.Second)

	dispatch(CreateDisconnectPlayerAction("X"))

	// Simulate player actions with dispatches
	time.Sleep(1 * time.Second) // Allow game loop to start

	dispatch(CreateDisconnectPlayerAction("O"))

	// Simulate player actions with dispatches
	time.Sleep(1 * time.Second) // Allow game loop to start
	// Assert equality using Testify
	assert.Equal(t, "X", state.winner, "Player 1 should win")

}
