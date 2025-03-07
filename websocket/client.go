package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/hexagun/common"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type     string          `json:"type"`
	GameID   string          `json:"gameId,omitempty"`
	PlayerID string          `json:"playerId,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}
		var msg common.IncomingMessage

		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("JSON Unmarshal error:", err)
			continue
		}

		c.Pool.Broadcast <- msg
		fmt.Printf("Message Received: %+v\n", message)
		fmt.Printf("Message Decoded: %+v\n", msg)
	}
}
