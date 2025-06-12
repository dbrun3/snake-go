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

func SnakeGame(game *game.GameState) {
	name := "dylan"
	color := objects.ColorBlue

	err := render.Init()
	if err != nil {
		panic(err)
	}
	defer render.Close()

	camera := render.CreateCamera()
	camera.SetSize(render.Size())

	var mySnake *objects.Snake = &objects.Snake{Dead: true, Body: []objects.Coord{{X: 0, Y: 0}}}
	var menu *render.Menu

	eventCh := make(chan termbox.Event, 4)
	lastSendTime := time.Now()

	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	tick := time.NewTicker(32 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case ev := <-eventCh:
			switch ev.Type {
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
					// Ensure menu exists and update it
					if menu == nil || !menu.Active {
						menu = render.NewMenu(&name, &color, camera)
					} else {
						menu.Update(ev)
					}
				} else {
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
						lastSendTime = now
					}
				}

			case termbox.EventResize:
				camera.SetSize(ev.Width, ev.Height)
			}

		case <-tick.C:

			// Handle snake respawn if menu just finished
			if mySnake.Dead && menu != nil && !menu.Active {
				newSnake := &objects.Snake{Id: game.ClientId, Name: name, Color: color}
				if game.IsServer() {
					game.AddSnake(newSnake)
				} else {
					data, _ := newSnake.Export()
					game.SendEvent("add_snake", data)
				}
				var err error
				mySnake, err = game.WaitForSnake(time.Second * 5)
				if err != nil {
					fmt.Println("Could not create snake", err)
					return
				}
				menu = nil // reset menu for next death
			}

			// Render everything else
			camera.FollowPos(mySnake.Head())
			render.Clear()
			render.RenderGameState(game, camera)
			if mySnake.Dead && menu != nil {
				menu.Draw()
			}
			render.Flush()
		}
	}
}
