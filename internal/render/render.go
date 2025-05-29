package render

import (
	"snake/internal/game"
	"snake/internal/objects"

	"github.com/nsf/termbox-go"
)

// Initialize window with resize callback
func Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	return nil
}

func Size() (int, int) {
	return termbox.Size()
}

// Close termbox
func Close() {
	termbox.Close()
}

// Clear the screen
func Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// Flush changes to screen
func Flush() {
	termbox.Flush()
}

// Render all snakes to the screen
func DrawSnakes(game *game.GameState, camera *Camera) {
	game.Mu.Lock()
	defer game.Mu.Unlock()
	for _, snake := range game.Snakes {
		if !snake.Dead {
			drawSnake(snake, camera)
		}
	}
}

func DrawFruits(game *game.GameState, camera *Camera) {
	game.Mu.Lock()
	defer game.Mu.Unlock()

	for x := range camera.cameraDim.X {
		for y := range camera.cameraDim.Y {
			fruit, fruitPresent := game.Fruits[camera.RealPos(objects.Coord{X: x, Y: y})]
			if fruitPresent {
				drawPoint(x, y, fruit.Color)
			}
		}
	}
}

// Draw a single point to screen (helper)
func drawPoint(x, y int, color objects.Color) {

	isTopOfCell := y%2 == 0

	cell := termbox.GetCell(x, y/2)
	existingCh := cell.Ch
	existingFg := cell.Fg

	if existingCh == ' ' {

		// Set new point, either top or bottom of full cell
		var newCh rune
		if isTopOfCell {
			newCh = '▀'
		} else {
			newCh = '▄'
		}

		termbox.SetCell(x, y/2, newCh, color.ToTermbox(), termbox.ColorDefault)

	} else if (existingCh == '▀') != isTopOfCell {

		if existingFg == color.ToTermbox() {

			// Replace with a full cell
			termbox.SetCell(x, y/2, '█', existingFg, termbox.ColorDefault)

		} else {

			// Overwrite background as the new cell
			termbox.SetCell(x, y/2, existingCh, existingFg, color.ToTermbox())
		}
	}
}

// Draw single snake (helper)
func drawSnake(snake *objects.Snake, camera *Camera) {
	for _, part := range snake.Body {
		canRender, renderPos := camera.RenderPos(part)
		if canRender {
			drawPoint(renderPos.X, renderPos.Y, snake.Color)
		}
	}
}
