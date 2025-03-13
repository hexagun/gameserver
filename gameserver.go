package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hexagun/common"
	"github.com/hexagun/gameserver/websocket"
)

var game *Game
var pool *websocket.Pool

type GameMessageDecoder struct {
}

func (g GameMessageDecoder) Decode(msg *common.IncomingMessage) {
	var action common.Action

	gameId, err := strconv.Atoi(msg.GameID)
	if err != nil {
		// ... handle error
		panic(err)
	}

	playerId, err := strconv.Atoi(msg.PlayerID)
	if err != nil {
		// ... handle error
		panic(err)
	}

	switch msg.Type {
	//case "gamestateupdate":
	//case "gameover":
	//case "error":
	//case "start":
	case "join":
		// string to int

		action = common.NewJoinAction(gameId, playerId)
	// case "leave":
	// 	action = common.NewLeaveAction(gameId, playerId)
	// case "move":
	// 	action = common.NewMoveAction(gameId, playerId)
	default:
		fmt.Println("Not Decoding stuff")
	}

	game.Dispatch(action)
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn:    conn,
		Pool:    pool,
		Decoder: GameMessageDecoder{},
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes(pool *websocket.Pool) {
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {

	port := ":8100"
	fmt.Println("Gameserver started on port %s", port)
	pool = websocket.NewPool()

	game = NewGame(pool)
	go game.GameLoop()
	setupRoutes(pool)
	log.Fatal(http.ListenAndServe(port, nil))
}
