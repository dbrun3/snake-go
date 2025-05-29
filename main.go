package main

import (
	"flag"
	"fmt"
	"snake/cmd"
	"snake/internal/client"
	"snake/internal/game"
	"snake/internal/objects"
	"snake/internal/server"
)

func main() {
	mode := flag.String("mode", "host", "host, client, server (default host)")
	port := flag.Int("port", 8080, "the port to open incoming connections (default 8080)")
	addr := flag.String("addr", "none", "the host/server to connect to (required if mode:client)")

	flag.Parse()
	if *mode != "server" && *mode != "host" && *mode != "client" {
		fmt.Println("Invalid mode")
		return
	}

	isServer := *mode != "client"
	isPlayer := *mode != "server"

	game := game.NewGameState(isServer)

	if isServer {
		go server.WsServer(game, *port)
	} else {
		go client.WsClient(game, *addr)
	}

	if isPlayer {
		go cmd.SnakeGame(game, "dylan", objects.ColorCyan)
	}

	game.GameLoop()
}
