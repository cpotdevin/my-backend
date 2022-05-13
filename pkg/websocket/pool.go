package websocket

import (
	"fmt"
	"log"
)

type Pool struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("%d", len(pool.Clients))})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("%d", len(pool.Clients))})
			}
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
