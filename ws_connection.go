package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsConnection(hub *Hub, w http.ResponseWriter, r *http.Request) {
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
}
