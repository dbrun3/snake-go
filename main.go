package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"snake/cmd"
	"snake/internal/client"
	"snake/internal/game"
	"snake/internal/server"
	"strings"
)

func main() {
	mode := flag.String("mode", "client", "host, client, server (default host)")
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
		openHostFile(addr)
		go client.WsClient(game, *addr)
	}

	fmt.Println("Snake-go!!!")

	if isPlayer {
		go cmd.SnakeGame(game)
	}

	game.GameLoop()
}

func openHostFile(addr *string) {
	if *addr == "none" {
		// open
		hf, err := os.Open(".snake.host")
		if err != nil {
			hf, err = os.Create(".snake.host")
			if err != nil {
				panic(err)
			}
		}

		// read
		buf := make([]byte, 48)
		n, err := hf.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		*addr = strings.TrimSpace(string(buf[:n]))

		// edit
		cmd.SelectHost(addr)

	}
	// write
	hf, err := os.Create(".snake.host")
	if err != nil {
		panic(err)
	}
	_, err = hf.Write([]byte(*addr))
	if err != nil {
		panic(err)
	}

}
