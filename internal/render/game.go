package render

import (
	"snake/internal/game"
	"snake/internal/objects"
	"sort"
	"strconv"

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

func RenderGameState(game *game.GameState, camera *Camera) {
	game.Mu.Lock()
	defer game.Mu.Unlock()

	drawFruits(game, camera)
	drawSnakes(game, camera)
	drawLeaderboard(game, camera)

}

func drawLeaderboard(game *game.GameState, camera *Camera) {

	s := 0
	sortedSnakes := make([]*objects.Snake, len(game.Snakes))
	for _, snake := range game.Snakes {
		sortedSnakes[s] = snake
		s++
	}

	sort.Slice(sortedSnakes, func(i, j int) bool {
		return sortedSnakes[i].Len > sortedSnakes[j].Len
	})

	x := camera.cameraDim.X - 15
	for y, snake := range sortedSnakes {
		line := snake.Name
		score := strconv.Itoa(snake.Len)

		buf := make([]byte, 15)
		copy(buf, line)
		copy(buf[(15-len(score)-1):], score)

		drawSentenceColor(x, y, string(buf), snake.Color.ToTermbox())
	}
}

// Render all snakes to the screen
func drawSnakes(game *game.GameState, camera *Camera) {
	for _, snake := range game.Snakes {
		if !snake.Dead {
			drawSnake(snake, camera)
		}
	}
}

func drawFruits(game *game.GameState, camera *Camera) {

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
