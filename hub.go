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
	connections []*websocket.Conn
	broadcast   chan Message
	unregister  chan *websocket.Conn
	register    chan *websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		connections: []*websocket.Conn{},
		broadcast:   make(chan Message),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
	}
}

func (h *Hub) startHub() {
	broadcast := func(message Message) {
		for _, conn := range h.connections {
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
			h.connections = append(h.connections, conn)
		case conn := <-h.unregister:
			fmt.Println("client disconnected")
			for i, c := range h.connections {
				if c == conn {
					h.connections = append(h.connections[:i], h.connections[i+1:]...)
					break
				}
			}
			broadcast(Message{0, "User disconnected", "Server Info"})
		}
	}
}
