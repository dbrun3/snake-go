package objects

import "fmt"

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewCord(x int, y int) Coord {
	return Coord{X: x, Y: y}
}

func Zero() Coord {
	return Coord{X: 0, Y: 0}
}

func (c *Coord) Equals(other Coord) bool {
	return c.X == other.X && c.Y == other.Y
}

func (c *Coord) Less(other Coord) bool {
	return c.X < other.X && c.Y < other.Y
}

func (c *Coord) LessOrEqual(other Coord) bool {
	return c.X <= other.X && c.Y <= other.Y
}

func (c *Coord) Greater(other Coord) bool {
	return c.X > other.X && c.Y > other.Y
}

func (c *Coord) GreaterOrEqual(other Coord) bool {
	return c.X >= other.X && c.Y >= other.Y
}

func (c *Coord) Add(other Coord) Coord {
	return Coord{X: c.X + other.X, Y: c.Y + other.Y}
}

func (c *Coord) Subtract(other Coord) Coord {
	return Coord{X: c.X - other.X, Y: c.Y - other.Y}
}

func (c *Coord) Translate(x int, y int) Coord {
	return Coord{X: c.X + x, Y: c.Y + y}
}

func (c Coord) MarshalText() (string, error) {
	return (fmt.Sprintf("%d,%d", c.X, c.Y)), nil
}

func (c *Coord) UnmarshalText(text string) error {
	_, err := fmt.Sscanf(text, "%d,%d", &c.X, &c.Y)
	return err
}
