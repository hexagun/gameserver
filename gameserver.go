package main

import (
	"encoding/json"
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
	var gameId int
	var playerId int

	if msg.GameID != "" {
		id, err := strconv.Atoi(msg.GameID)
		if err != nil {
			// ... handle error
			panic(err)
		}
		gameId = id
	}

	if msg.PlayerID != "" {
		id, err := strconv.Atoi(msg.PlayerID)
		if err != nil {
			// ... handle error
			panic(err)
		}
		playerId = id
	}

	switch msg.Type {
	//case "gamestateupdate":
	//case "gameover":
	//case "error":
	case "start":
		//action = common.NewStartAction(gameId, playerId)
	case "join":
		action = common.NewJoinAction(gameId, playerId)
	// case "leave":
	// 	action = common.NewLeaveAction(gameId, playerId)
	case "move":
		var movePayload struct {
			Row int `json:"row"`
			Col int `json:"col"`
		}
		if err := json.Unmarshal(msg.Payload, &movePayload); err == nil {
			fmt.Printf("Player made a move at row %d, col %d\n", movePayload.Row, movePayload.Col)
			action = common.NewPlayerMoveAction(gameId, playerId, common.PlayerMovePayload{
				Row: movePayload.Row,
				Col: movePayload.Col,
			})
		}

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
	idUser := r.URL.Query().Get("id")
	client := &websocket.Client{
		ID:      idUser,
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
