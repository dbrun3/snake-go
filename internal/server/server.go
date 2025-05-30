package server

import (
	"snake/internal/events"
	"snake/internal/game"
	"snake/internal/objects"
)

type Server struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	events chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Game state
	game *game.GameState
}

type Message struct {
	message []byte
	sender  string
}

func newGameServer(game *game.GameState) *Server {
	return &Server{
		events:     make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		game:       game,
	}
}

func (s *Server) run() {
	go s.eventServer()
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true

			data, _ := s.game.Export(client.id)
			event := events.NewEvent("init", data)
			e, _ := events.MarshalEvent(event)
			client.send <- e

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}

			snake := &objects.Snake{Id: client.id}
			s.game.RemoveSnake(snake)

		case event := <-s.events:
			s.game.HandleEvent(event.sender, event.message)
		}
	}
}

func (s *Server) eventServer() {
	for event := range s.game.Send {
		for client := range s.clients {
			select {
			case client.send <- event:
			default:
				close(client.send)
				delete(s.clients, client)
			}
		}
	}
}
