package render

import "snake/internal/objects"

type Camera struct {
	offsetPos objects.Coord
	cameraDim objects.Coord
}

func CreateCamera() *Camera {
	return &Camera{}
}

func (c *Camera) FollowPos(pos objects.Coord) {
	targetPos := objects.NewCord(
		pos.X-(c.cameraDim.X/2),
		pos.Y-(c.cameraDim.Y/2),
	)

	// Calculate direction (normalized to ±1)
	dx := targetPos.X - c.offsetPos.X
	dy := targetPos.Y - c.offsetPos.Y

	// Move at least 1px toward the target (even if diff is small)
	if dx != 0 {
		c.offsetPos.X += dx / abs(dx) // dx/|dx| = ±1
	}
	if dy != 0 {
		c.offsetPos.Y += dy / abs(dy) // dy/|dy| = ±1
	}
}

// Helper for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
