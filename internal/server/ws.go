package server

import (
	"fmt"
	"net/http"
	"snake/internal/game"
)

func WsServer(game *game.GameState, port int) error {
	server := newGameServer(game)
	go server.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return err
}
