package game

import (
	"encoding/json"
	"fmt"
	"snake/internal/objects"
)

type ExportState struct {
	Snakes   map[string]objects.Snake `json:"snakes"`
	Fruits   map[string]objects.Fruit `json:"fruits"` // "x,y" -> Fruit
	ClientId string                   `json:"id"`
}

func (gs *GameState) Import(data []byte) {
	var eS ExportState
	if err := json.Unmarshal(data, &eS); err != nil {
		fmt.Printf("JSON unmarshal failed: %s\n", err)
	}

	game := importFrom(&eS)

	gs.Mu.Lock()
	gs.Snakes = game.Snakes
	gs.Fruits = game.Fruits
	gs.ClientId = game.ClientId

	gs.Mu.Unlock()
}

func (gs *GameState) Export(clientId string) ([]byte, error) {
	exportState := gs.exportTo()
	exportState.ClientId = clientId
	return json.Marshal(exportState)
}

func importFrom(exportState *ExportState) *GameState {

	gameState := &GameState{
		Snakes:   make(map[string]*objects.Snake),
		Fruits:   make(map[objects.Coord]objects.Fruit),
		ClientId: exportState.ClientId,
	}

	for name, snake := range exportState.Snakes {
		snakeCopy := snake // local copy prevents aliasing ig
		gameState.Snakes[name] = &snakeCopy
	}

	// unmarshal cords to ints
	for coordStr, fruit := range exportState.Fruits {
		var coord objects.Coord
		if err := coord.UnmarshalText(coordStr); err == nil {
			gameState.Fruits[coord] = fruit
		}
	}

	return gameState
}

func (gs *GameState) exportTo() ExportState {
	gs.Mu.RLock() // Read lock for thread safety
	defer gs.Mu.RUnlock()

	// deep copy
	Snakes := make(map[string]objects.Snake, len(gs.Snakes))
	for id, snakePtr := range gs.Snakes {
		Snakes[id] = *snakePtr
	}

	Fruits := make(map[string]objects.Fruit, len(gs.Fruits))
	for coord, fruit := range gs.Fruits {
		key, _ := coord.MarshalText() // converts Coord to string
		Fruits[key] = fruit
	}

	return ExportState{
		Snakes: Snakes,
		Fruits: Fruits,
	}
}
