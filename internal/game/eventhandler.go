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
	case "kill_snake":
		snake, err := objects.ImportSnake(event.Data)
		if err != nil {
			fmt.Print(err)
			return
		}
		gs.KillSnake(snake)

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
	}
}

func (gs *GameState) UpdateSnake(updatedSnake *objects.Snake) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	id := updatedSnake.Id

	snake, exists := gs.Snakes[id]
	if !exists || gs.ClientId == id {
		return
	}

	// Update direction and other states
	snake.ChangeDir(updatedSnake.Dir)
	snake.Dead = updatedSnake.Dead
	snake.Len = updatedSnake.Len
	snake.Speed = updatedSnake.Speed

	// Set a target to handle desync
	snake.Target = updatedSnake.Head()

	if gs.IsServer() {
		// Rebroadcast update
		data, _ := updatedSnake.Export()
		gs.SendEvent("update_snake", data)
		return
	}
}

func (gs *GameState) KillSnake(updatedSnake *objects.Snake) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	id := updatedSnake.Id

	snake, exists := gs.Snakes[id]
	if !exists {
		return
	}

	snake.Dead = true

	r := rand.New(rand.NewSource(int64(updatedSnake.Len)))
	for index, part := range snake.Body {
		if index%2 == 0 {
			objects.CreateFruit(
				snake.Color,
				objects.Coord{X: part.X + r.Intn(2) - 1, Y: part.Y + r.Intn(2) - 1},
				&gs.Fruits)

		}
	}

	// no server rebroadcast necessary, server-only event
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
