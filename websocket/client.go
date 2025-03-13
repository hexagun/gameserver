package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/hexagun/common"
)

//type MessageDecoder func(message *common.IncomingMessage)

type MessageDecoder interface {
	Decode(message *common.IncomingMessage)
}

type Client struct {
	ID      string
	Conn    *websocket.Conn
	Pool    *Pool
	Decoder MessageDecoder
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()

		fmt.Printf("Message Received: %+v\n", message)
		if err != nil {
			log.Println(err)
			return
		}
		var msg common.IncomingMessage

		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("JSON Unmarshal error:", err)
			continue
		}

		fmt.Printf("Message Decoded: %+v\n", msg)

		c.DecodeMessage(&msg)

	}
}

func (c *Client) DecodeMessage(msg *common.IncomingMessage) {
	c.Decoder.Decode(msg)
}
