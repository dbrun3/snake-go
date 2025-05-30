package game

import (
	"math/rand"
	"os"
	"os/signal"
	"snake/internal/objects"
	"time"
)

const FRUIT_PER_TICK = 20
const MAX_FRUIT = 200
const MAP_SIZE = 300
const TICK_DURATION = 50

func (gs *GameState) GameLoop() {

	// initialize timers and random seeds
	tick := time.NewTicker(TICK_DURATION * time.Millisecond)
	serverRand := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	// only the server creates seeds
	if gs.IsServer() {
		defer serverRand.Stop()
		gs.r = rand.New(rand.NewSource(time.Now().UnixMilli()))
	} else {
		serverRand.Stop()
	}

	// handle interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {

		select {

		case <-interrupt: // Escape
			return

		case <-serverRand.C: // Server will generate new seeds every second
			gs.PlantSeed(int(time.Now().UnixMilli()))

		case <-tick.C: // Regular game tick 50 milliseconds

			gs.Mu.Lock()
			for _, snake := range gs.Snakes {

				// Dont update dead snakes
				if snake.Dead {
					continue
				}

				// Move
				snake.Move()

				// Update head map
				gs.heads[snake.Head()] = snake

				// Eat fruit
				_, fruitPresent := gs.Fruits[snake.Head()]
				if fruitPresent {
					delete(gs.Fruits, snake.Head())
					snake.Eat()
				}

				// Another snake's head crashes into body
				for _, b := range snake.Body {
					hitSnake, hit := gs.heads[b]

					if hitSnake != snake && hit && !hitSnake.Dead {

						hitSnake.Dead = true

						r := rand.New(rand.NewSource(int64(hitSnake.Len)))
						for index, part := range hitSnake.Body {
							if index%2 == 0 {
								objects.CreateFruit(
									hitSnake.Color,
									objects.Coord{X: part.X + r.Intn(2) - 1, Y: part.Y + r.Intn(2) - 1},
									&gs.Fruits)

							}
						}
					}
				}
			}
			gs.Mu.Unlock()
		}

	}
}
