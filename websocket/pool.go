package websocket

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/hexagun/common"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Broadcast  chan common.OutgoingMessage
	Send       chan common.OutgoingMessage
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan common.OutgoingMessage),
		Send:       make(chan common.OutgoingMessage),
	}
}

func GenerateCustomID() int {
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return int(randomNumber.Int64())
}

func (pool *Pool) Start() {
	gameId := GenerateCustomID()
	for {
		select {
		case client := <-pool.Register:

			msg := &common.IncomingMessage{
				Type:     "join",
				GameID:   fmt.Sprintf("%d", gameId),
				PlayerID: client.ID,
				Payload:  nil,
			}
			client.Decoder.Decode(msg)
			pool.Clients[client.ID] = client
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// 	fmt.Println(client)

			// 	//client.Conn.WriteJSON(common.OutgoingMessage{Type: "Join"})
			// }
			break
		case client := <-pool.Unregister:

			msg := &common.IncomingMessage{
				Type:     "leave",
				GameID:   fmt.Sprintf("%d", gameId),
				PlayerID: client.ID,
				Payload:  nil,
			}
			client.Decoder.Decode(msg)

			delete(pool.Clients, client.ID)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// 	client.Conn.WriteJSON(common.OutgoingMessage{Type: ""})
			// }
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for _, client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			break
		case message := <-pool.Send:

			for id, client := range pool.Clients {
				if id == message.PlayerID {
					fmt.Println("Sending message to client:", id)
					fmt.Println("Sending message:", message)
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
					}
				}
			}
			fmt.Println("Message Sent")
			break
		}
	}
}
