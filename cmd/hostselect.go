package cmd

import (
	"fmt"
	"os"
	"snake/internal/render"
	"time"

	"github.com/nsf/termbox-go"
)

func SelectHost(host *string) {

	err := render.Init()
	if err != nil {
		panic(err)
	}
	defer render.Close()

	eventCh := make(chan termbox.Event, 4)
	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	camera := render.CreateCamera()
	camera.SetSize(render.Size())

	hostSelection := &render.HostSelection{
		Active: true,
		Host:   host,
		Camera: camera,
	}

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
				hostSelection.Update(ev)
			case termbox.EventResize:
				camera.SetSize(ev.Width, ev.Height)
			}
		case <-tick.C:
			render.Clear()
			hostSelection.Draw()
			if !hostSelection.Active {
				return
			}
			render.Flush()
		}
	}
}
