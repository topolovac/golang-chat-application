package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("public")))

	// Start hub and handle websocket connections
	hub := newHub()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsConnection(hub, w, r)
	})
	go hub.startHub()

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server listening on port %s...", server.Addr)
	log.Fatal(server.ListenAndServe())
}
