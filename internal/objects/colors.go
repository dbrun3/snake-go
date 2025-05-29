package objects

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

type Color termbox.Attribute

const (
	ColorRed     Color = Color(termbox.ColorRed)
	ColorGreen   Color = Color(termbox.ColorGreen)
	ColorYellow  Color = Color(termbox.ColorYellow)
	ColorBlue    Color = Color(termbox.ColorBlue)
	ColorMagenta Color = Color(termbox.ColorMagenta)
	ColorCyan    Color = Color(termbox.ColorCyan)
	ColorWhite   Color = Color(termbox.ColorWhite)
)

var AllColors = []Color{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorCyan,
	ColorWhite,
}

var ColorNames = []string{"Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White"}

// ToTermbox converts Color to termbox.Attribute
func (c Color) ToTermbox() termbox.Attribute {
	return termbox.Attribute(c)
}

// RandomColor returns a random color from AllColors
func RandomColor() Color {
	return AllColors[rand.Intn(len(AllColors))]
}
