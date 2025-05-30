package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"snake/internal/events"
	"snake/internal/objects"
)

func (gs *GameState) SendEvent(name string, data []byte) {
	event := events.NewEvent(name, data)
	e, _ := events.MarshalEvent(event)
	gs.Send <- e
}

func (gs *GameState) HandleEvent(sender string, e []byte) {

	event, _ := events.UnmarshalEvent(e)

	switch event.Type {

	case "init":
		gs.Import(event.Data)

	case "add_snake":
		snake, err := objects.ImportSnake(event.Data)
		if err != nil || (gs.IsServer() && snake.Id != sender) {
			return
		}
		gs.AddSnake(snake)

	case "update_snake":
		snake, err := objects.ImportSnake(event.Data)
		if err != nil || (gs.IsServer() && snake.Id != sender) {
			return
		}
		gs.UpdateSnake(snake)

	case "remove_snake":
		snake, err := objects.ImportSnake(event.Data)
		if err != nil {
			fmt.Print(err)
			return
		}
		gs.RemoveSnake(snake)

	case "plant_seed":
		var seed Seed
		json.Unmarshal(event.Data, &seed)
		gs.PlantSeed(seed.Seed)
	}
}

func (gs *GameState) AddSnake(snake *objects.Snake) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	var x, y int
	if gs.IsServer() {

		x = 0 //gs.r.Intn(MAP_SIZE/2) - (MAP_SIZE / 4)
		y = 0 //gs.r.Intn(MAP_SIZE/2) - (MAP_SIZE / 4)
	} else {
		x = snake.Head().X
		y = snake.Head().Y
	}

	// Find existing or allocate new snake
	newSnake, found := gs.Snakes[snake.Id]
	if !found {
		newSnake = objects.CreateSnake(snake.Id)
	}

	newSnake.Name = snake.Name
	newSnake.Color = snake.Color
	clear(newSnake.Body)
	newSnake.Body = append(newSnake.Body, objects.NewCord(x, y))
	newSnake.Len = 1
	newSnake.Dead = false

	gs.Snakes[snake.Id] = newSnake

	if gs.IsServer() {
		// Rebroadcast new snake
		data, _ := newSnake.Export()
		gs.SendEvent("add_snake", data)
		return
	}
}

func (gs *GameState) UpdateSnake(updatedSnake *objects.Snake) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	id := updatedSnake.Id

	snake, exists := gs.Snakes[id]
	if !exists {
		return
	}

	// Update direction
	snake.ChangeDir(updatedSnake.Dir)

	if gs.IsServer() {
		// Rebroadcast update
		data, _ := snake.Export()
		gs.SendEvent("update_snake", data)
		return
	}

	// Sync additional states with server
	snake.Dead = updatedSnake.Dead
	snake.Len = updatedSnake.Len

	updatedLen := len(updatedSnake.Body)
	snakeLen := len(snake.Body)

	// If the updated body is longer than the current body, just replace entirely
	if updatedLen >= snakeLen {
		snake.Body = updatedSnake.Body
		return
	}

	h := snake.Head()
	for i := len(updatedSnake.Body) - 1; i >= 0; i-- {
		if h.Equals(updatedSnake.Body[i]) {
			snake.Body = append(snake.Body, updatedSnake.Body[i+1:]...)
			return
		}
	}

	// fallback
	snake.Body = append(snake.Body[:snakeLen-updatedLen], updatedSnake.Body...)

}

func (gs *GameState) RemoveSnake(snake *objects.Snake) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if gs.Snakes[snake.Id] != nil {
		delete(gs.heads, gs.Snakes[snake.Id].Head())
		delete(gs.Snakes, snake.Id)
	}

	if gs.IsServer() {
		data, _ := snake.Export()
		gs.SendEvent("remove_snake", data)
		return
	}
}

type Seed struct {
	Seed int `json:"seed"`
}

func (gs *GameState) PlantSeed(seed int) {
	// Seed random events

	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	gs.r = rand.New(rand.NewSource(int64(seed)))
	if len(gs.Fruits) < MAX_FRUIT {
		for range FRUIT_PER_TICK {
			objects.CreateFruit(
				objects.RandomColor(),
				objects.Coord{X: gs.r.Intn(MAP_SIZE) - (MAP_SIZE / 2), Y: gs.r.Intn(MAP_SIZE) - (MAP_SIZE / 2)},
				&gs.Fruits)
		}
	}

	if gs.IsServer() {
		data, _ := json.Marshal(Seed{Seed: seed})
		gs.SendEvent("plant_seed", data)
		return
	}
}
