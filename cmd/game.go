package cmd

import (
	"fmt"
	"os"
	"snake/internal/game"
	"snake/internal/objects"
	"snake/internal/render"
	"time"

	"github.com/nsf/termbox-go"
)

// Todo refactor to handle inputs in separate function from rendering. add death screen menu
func SnakeGame(game *game.GameState, name string, color objects.Color) {

	// initialize window
	err := render.Init()
	if err != nil {
		panic(err)
	}
	defer render.Close()

	// start camera w a dead snake to watch
	camera := render.CreateCamera()
	camera.SetSize(render.Size())
	var mySnake *objects.Snake = &objects.Snake{Dead: true, Body: []objects.Coord{{X: 0, Y: 0}}}

	// Event channel
	eventCh := make(chan termbox.Event, 4)
	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	// Event and render loop
	tick := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer tick.Stop()
	for {

		select {

		// event loop
		case ev := <-eventCh:
			switch ev.Type {

			// key press
			case termbox.EventKey:

				if ev.Key == termbox.KeyEsc {
					render.Close()
					process, err := os.FindProcess(os.Getpid())
					if err != nil {
						fmt.Println(err)
					}
					process.Signal(os.Interrupt)
					return
				}

				if mySnake.Dead {

					// on start/death keypress to submit a new snake
					newSnake := &objects.Snake{Id: game.ClientId, Name: name, Color: color}

					if game.IsServer() {
						game.AddSnake(newSnake)
					} else {
						data, _ := newSnake.Export()
						game.SendEvent("add_snake", data)
					}

					mySnake, err = game.WaitForSnake(time.Second * 5)
					if err != nil {
						fmt.Println("Could not create snake", err)
						return
					}

				} else {

					switch ev.Key {

					case termbox.KeyArrowRight:
						mySnake.ChangeDir(objects.RIGHT)

					case termbox.KeyArrowLeft:
						mySnake.ChangeDir(objects.LEFT)

					case termbox.KeyArrowUp:
						mySnake.ChangeDir(objects.UP)

					case termbox.KeyArrowDown:
						mySnake.ChangeDir(objects.DOWN)
					}

					if game.IsServer() {
						game.UpdateSnake(mySnake)
					} else {
						data, _ := mySnake.Export()
						game.SendEvent("update_snake", data)
					}
				}

			// resize event, update camera
			case termbox.EventResize:
				camera.SetSize(ev.Width, ev.Height)
			}

		// render loop
		case <-tick.C:
			// camera pointed at snake head
			camera.SetPos(mySnake.Head())

			// main render loop
			render.Clear()
			render.DrawFruits(game, camera)
			render.DrawSnakes(game, camera)
			render.Flush()
		}
	}
}
