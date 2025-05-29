package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"snake/internal/game"
	"time"

	"github.com/gorilla/websocket"
)

func WsClient(game *game.GameState, address string) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: address, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}

			reader := bytes.NewReader(message)
			decoder := json.NewDecoder(reader)

			// Keep decoding until we've processed all JSON objects in the message
			for {
				var data json.RawMessage
				err := decoder.Decode(&data)
				if err != nil {
					if err == io.EOF {
						break // We've processed all complete JSON objects
					}
					fmt.Println("Error decoding JSON:", err)
					break
				}
				game.HandleEvent("server", data)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case t := <-game.Send:
			err := c.WriteMessage(websocket.TextMessage, t)
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		case <-interrupt:
			fmt.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
