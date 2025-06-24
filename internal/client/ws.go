package client

import (
	"bytes"
	"encoding/json"
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
		panic(err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				panic(err)
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
					panic(err)
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
				panic(err)
			}
		case <-interrupt:

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				panic(err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
