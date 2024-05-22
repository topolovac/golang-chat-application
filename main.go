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
	}

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		hub.connections = append(hub.connections, conn)

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
				log.Fatal(err)
			}

			fmt.Printf("new message: %v\n", message)
		}
	})

	log.Printf("Server listening on port %s...", server.Addr)
	log.Fatal(server.ListenAndServe())
}
