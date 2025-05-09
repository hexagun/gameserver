package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hexagun/common"
	"github.com/hexagun/gameserver/websocket"
)

var game *Game
var pool *websocket.Pool

type GameMessageDecoder struct {
	messageDecoder common.MessageDecoder
}

func (g GameMessageDecoder) Decode(msg *common.IncomingMessage) {

	action := g.messageDecoder.Decode(msg)
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
