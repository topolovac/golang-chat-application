package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Message struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	User string `json:"user"`
}

type Hub struct {
	connections map[*websocket.Conn]bool
	broadcast   chan Message
	unregister  chan *websocket.Conn
	register    chan *websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		connections: make(map[*websocket.Conn]bool),
		broadcast:   make(chan Message),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
	}
}

func (h *Hub) startHub() {
	broadcast := func(message Message) {
		for conn := range h.connections {
			if err := conn.WriteJSON(message); err != nil {
				fmt.Printf("error writing to connection: %v\n", err)
			}
		}
	}
	for {
		select {
		case message := <-h.broadcast:
			fmt.Printf("broadcasting message: %v\n", message)
			broadcast(message)
		case conn := <-h.register:
			fmt.Println("new connection")
			broadcast(Message{0, "User connected", "Server Info"})
			h.connections[conn] = true
		case conn := <-h.unregister:
			fmt.Println("client disconnected")
			delete(h.connections, conn)
			broadcast(Message{0, "User disconnected", "Server Info"})
		}
	}
}
