package render

import (
	"snake/internal/objects"
	"strings"

	"github.com/nsf/termbox-go"
)

type Menu struct {
	Active           bool
	state            int // 0 = name input, 1 = color selection
	Colors           []objects.Color
	PromptX, PromptY int
	InputBoxY        int
	Camera           *Camera
	NamePtr          *string
	ColorPtr         *int
}

type HostSelection struct {
	Active bool
	Host   *string
	Camera *Camera
}

func NewMenu(name *string, color *int, camera *Camera) *Menu {
	width := camera.cameraDim.X
	height := camera.cameraDim.Y

	menuWidth := width / 2
	menuHeight := height / 2
	startX := (width - menuWidth) / 2
	startY := (height - menuHeight) / 2

	return &Menu{
		Active:    true,
		state:     0,
		Colors:    objects.AllColors,
		PromptX:   startX + 2,
		PromptY:   startY + 1,
		InputBoxY: startY + 3,
		Camera:    camera,
		NamePtr:   name,
		ColorPtr:  color,
	}
}

func (m *Menu) Update(event termbox.Event) {
	if !m.Active {
		return
	}
	if event.Type == termbox.EventKey {
		switch m.state {
		case 0:
			if typeEvent(m.NamePtr, event) {
				m.state = 1
			}
		case 1:
			if horizontalSelectEvent(m.ColorPtr, len(m.Colors), event) {
				m.Active = false
			}
		}
	}
}

func (h *HostSelection) Update(event termbox.Event) {
	if !h.Active {
		return
	}
	if event.Type == termbox.EventKey {
		if typeEvent(h.Host, event) {
			h.Active = false
		}
	}
}

func (m *Menu) Draw() {
	if !m.Active {
		return
	}

	width := m.Camera.cameraDim.X
	height := m.Camera.cameraDim.Y
	menuHeight := 30
	menuWidth := 50

	if width < menuWidth || height < menuHeight {
		return // Too small to display menu meaningfully
	}

	vOffset := (height - menuHeight) / 4 // Center the entire menu vertically

	// Draw background rectangle for the entire menu
	startX := (width - menuWidth) / 2
	startY := vOffset
	clearRectangle(startX, startY, menuWidth, menuHeight)

	// Center prompt
	promptText := "Enter your name:"
	promptX := (width - len(promptText)) / 2
	drawSentence(promptX, startY+1, promptText)

	// Draw name input box centered
	nameBoxWidth := max(len(*m.NamePtr), 10)
	nameBoxX := (width - (nameBoxWidth + 2)) / 2
	nameBoxY := startY + 3
	drawOutline(nameBoxX, nameBoxY, nameBoxWidth+2, 3)
	drawSentence(nameBoxX+1, nameBoxY+1, padRight(*m.NamePtr, 10))

	// Draw color selection after name selection
	if m.state > 0 {

		colorPromptY := nameBoxY + 5
		promptText := "Select a color with the arrows, followed by enter:"
		promptX := (width - len(promptText)) / 2
		drawSentence(promptX, colorPromptY, promptText)

		colorY := colorPromptY + 3
		totalColorWidth := len(m.Colors) * 6
		colorStartX := (width - totalColorWidth) / 2
		for i, col := range m.Colors {
			x := colorStartX + i*6
			drawColorBox(x, colorY, col, i == int(*m.ColorPtr))
		}
	}
}

func (h *HostSelection) Draw() {
	if !h.Active {
		return
	}

	width := h.Camera.cameraDim.X
	height := h.Camera.cameraDim.Y
	menuHeight := 15
	menuWidth := 50

	if width < menuWidth || height < menuHeight {
		return // Too small to display menu meaningfully
	}

	vOffset := (height - menuHeight) / 4 // Center the entire menu vertically

	// Draw background rectangle for the entire menu
	startX := (width - menuWidth) / 2
	startY := vOffset
	clearRectangle(startX, startY, menuWidth, menuHeight)

	// Center prompt
	promptText := "Enter a host to connect to"
	promptX := (width - len(promptText)) / 2
	drawSentence(promptX, startY+1, promptText)
	startY += 1

	// Draw name input box centered
	nameBoxWidth := max(len(*h.Host), 20)
	nameBoxX := (width - (nameBoxWidth + 2)) / 2
	nameBoxY := startY + 3
	drawOutline(nameBoxX, nameBoxY, nameBoxWidth+2, 3)
	drawSentence(nameBoxX+1, nameBoxY+1, padRight(*h.Host, 10))
}

func typeEvent(text *string, event termbox.Event) bool {
	switch event.Key {
	case termbox.KeyEnter:
		if len(*text) > 0 {
			return true
		}
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if len(*text) > 0 {
			*text = (*text)[:len(*text)-1]
		}
	default:
		if event.Ch != 0 {
			*text += string(event.Ch)
		}
	}
	return false
}

func horizontalSelectEvent(selection *int, max int, event termbox.Event) bool {
	switch event.Key {
	case termbox.KeyArrowLeft:
		if *selection > 0 {
			*selection--
		}
	case termbox.KeyArrowRight:
		if int(*selection) < max-1 {
			*selection++
		}
	case termbox.KeyEnter:
		return true
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func padRight(str string, minLen int) string {
	if len(str) >= minLen {
		return str
	}
	return str + strings.Repeat(" ", minLen-len(str))
}

func drawOutline(x, y, w, h int) {
	drawChar(x, y, '┌', termbox.ColorWhite)
	drawChar(x+w-1, y, '┐', termbox.ColorWhite)
	drawChar(x, y+h-1, '└', termbox.ColorWhite)
	drawChar(x+w-1, y+h-1, '┘', termbox.ColorWhite)

	for i := range w - 2 {
		drawChar(x+1+i, y, '─', termbox.ColorWhite)
		drawChar(x+1+i, y+h-1, '─', termbox.ColorWhite)
	}
	for i := range h - 2 {
		drawChar(x, y+1+i, '|', termbox.ColorWhite)
		drawChar(x+w-1, y+1+i, '|', termbox.ColorWhite)
	}
}

func drawColorBox(x, y int, color objects.Color, selected bool) {
	if selected {
		drawOutline(x, y, 4, 3)
	}
	drawChar(x+1, y+1, '█', color.ToTermbox())
	drawChar(x+2, y+1, '█', color.ToTermbox())
}

func clearRectangle(x, y, w, h int) {
	for i := range w {
		for j := range h {
			termbox.SetCell(x+i, y+j, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func drawSentence(x, y int, sentence string) {
	for i, c := range sentence {
		drawChar(x+i, y, c, termbox.ColorWhite)
	}
}

func drawChar(x, y int, char rune, color termbox.Attribute) {
	termbox.SetCell(x, y, char, color, termbox.ColorDefault)
}
