package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func broadcast(hub *Hub, message Message) {
	for _, conn := range hub.connections {
		if err := conn.WriteJSON(message); err != nil {
			fmt.Printf("error writing to connection: %v\n", err)
		}
	}
}

func (h *Hub) startHub() {
	for {
		select {
		case message := <-h.broadcast:
			fmt.Printf("broadcasting message: %v\n", message)
			broadcast(h, message)
		case conn := <-h.register:
			fmt.Println("new connection")
			broadcast(h, Message{0, "User connected", "Server Info"})
			h.connections = append(h.connections, conn)
		case conn := <-h.unregister:
			fmt.Println("client disconnected")
			for i, c := range h.connections {
				if c == conn {
					h.connections = append(h.connections[:i], h.connections[i+1:]...)
					break
				}
			}
			broadcast(h, Message{0, "User disconnected", "Server Info"})
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("public")))

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	hub := &Hub{
		connections: []*websocket.Conn{},
		broadcast:   make(chan Message),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
	}

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		hub.register <- conn

		err = conn.WriteJSON(&Message{
			Id:   0,
			Text: "Welcome to the chat!",
			User: "Server Info",
		})
		if err != nil {
			log.Fatal(err)
		}

		for {
			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				fmt.Println("error reading messages")
				conn.Close()
				hub.unregister <- conn
				break
			}

			hub.broadcast <- message
		}
	})

	go hub.startHub()

	log.Printf("Server listening on port %s...", server.Addr)
	log.Fatal(server.ListenAndServe())
}
