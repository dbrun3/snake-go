package render

import "snake/internal/objects"

type Camera struct {
	offsetPos objects.Coord
	cameraDim objects.Coord
}

func CreateCamera() *Camera {
	return &Camera{}
}

func (c *Camera) SetPos(pos objects.Coord) {
	c.offsetPos = objects.NewCord(pos.X-(c.cameraDim.X/2), pos.Y-(c.cameraDim.Y/2))
}

// Sets the camera dimensions given terminal cols and rows
func (c *Camera) SetSize(cols int, rows int) {
	c.cameraDim = objects.NewCord(cols, rows*2)
}

// Returns canRender? and the resulting position
func (c *Camera) RenderPos(pos objects.Coord) (bool, objects.Coord) {
	newPos := pos.Subtract(c.offsetPos)
	canRender := newPos.Less(c.cameraDim) && newPos.Greater(objects.Zero())
	return canRender, newPos
}

func (c *Camera) RealPos(pos objects.Coord) objects.Coord {
	newPos := pos.Add(c.offsetPos)
	return newPos
}
