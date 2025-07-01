package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
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

var secretKey = []byte("secret-key")

func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	tokenStr := r.URL.Query().Get("token")
	//idUser := r.URL.Query().Get("id")

	token, errToken := validateToken(tokenStr)
	if errToken != nil {
		fmt.Println("Invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Optional: extract claims
	// Should probably be id and a username retrieval from other service
	//username := ""
	id := 0
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("Token is valid!")
		//username = string(claims["username"])
		fmt.Println("Name:", claims["username"])
		fmt.Println("Id:", claims["uid"])
		fmt.Println("Expires:", claims["exp"])
		//username = claims["username"].(string)
		id = int(claims["uid"].(float64))
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		ID:      strconv.Itoa(id),
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
