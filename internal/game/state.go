package game

import (
	"context"
	"fmt"
	"math/rand"
	"snake/internal/objects"
	"sync"
	"time"
)

type GameState struct {
	// data
	Snakes   map[string]*objects.Snake
	Fruits   map[objects.Coord]objects.Fruit
	ClientId string

	// sync
	Send chan []byte
	Mu   sync.RWMutex

	// hidden
	heads map[objects.Coord]*objects.Snake
	r     *rand.Rand
}

func NewGameState(isServer bool) *GameState {

	state := &GameState{
		Snakes: make(map[string]*objects.Snake),       // Required
		Fruits: make(map[objects.Coord]objects.Fruit), // Required
		Send:   make(chan []byte),

		heads: make(map[objects.Coord]*objects.Snake),
	}

	if !isServer {
		state.ClientId = "unintialized"
	} else {
		state.ClientId = "server"
	}

	return state
}

func (gs *GameState) WaitForSnake(timeout time.Duration) (*objects.Snake, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for snake... ID: %s", gs.ClientId)
		case <-ticker.C:
			if snake, found := gs.getSnake(gs.ClientId); found {
				return snake, nil
			}
		}
	}
}

func (gs *GameState) getSnake(id string) (*objects.Snake, bool) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	snake, found := gs.Snakes[id]
	return snake, found
}

func (gs *GameState) IsServer() bool {
	return gs.ClientId == "server"
}
