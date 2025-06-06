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

// TODO: Add init/death screen menu to set name and color, rather than using defaults
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
	lastSendTime := time.Now() // the input throttle timer
	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	// Event and render loop
	tick := time.NewTicker(32 * time.Millisecond)
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

					// snake must first be registered on the server if not host
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
					// Throttle game inputs to prevent congestion (websockets will block to process everything in correct order)
					now := time.Now()
					if now.Sub(lastSendTime) > 100*time.Millisecond {
						switch ev.Key {

						case termbox.KeyArrowRight:
							mySnake.ChangeDir(objects.RIGHT)

						case termbox.KeyArrowLeft:
							mySnake.ChangeDir(objects.LEFT)

						case termbox.KeyArrowUp:
							mySnake.ChangeDir(objects.UP)

						case termbox.KeyArrowDown:
							mySnake.ChangeDir(objects.DOWN)

						case termbox.KeySpace:
							mySnake.ChangeSpeed()
						}
						data, _ := mySnake.Export()
						game.SendEvent("update_snake", data)
					}
					lastSendTime = now
				}

			// resize event, update camera
			case termbox.EventResize:
				camera.SetSize(ev.Width, ev.Height)
			}

		// render loop
		case <-tick.C:
			// camera pointed at snake head
			camera.FollowPos(mySnake.Head())

			// main render loop
			render.Clear()
			render.RenderGameState(game, camera)
			render.Flush()
		}
	}
}
