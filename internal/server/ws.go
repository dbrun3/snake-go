package server

import (
	"fmt"
	"log"
	"net/http"
	"snake/internal/game"
)

func WsServer(game *game.GameState, port int) error {
	server := newGameServer(game)
	go server.run()

	// WebSocket endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP server on %s", addr)

	// Block here until the server exits
	return http.ListenAndServe(addr, nil)
}
