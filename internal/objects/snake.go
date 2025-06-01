package objects

import (
	"encoding/json"
	"math"
)

const (
	UP int = iota
	DOWN
	LEFT
	RIGHT
)

type Snake struct {
	Dir   int     `json:"dir"`
	Len   int     `json:"len"`
	Color Color   `json:"color"`
	Name  string  `json:"name"`
	Id    string  `json:"id"`
	Body  []Coord `json:"body"`
	Dead  bool    `json:"dead"`
	Speed bool    `json:"speed"`

	Target Coord
}

func CreateSnake(id string) *Snake {
	return &Snake{Dir: RIGHT, Len: 1, Name: "", Id: id, Color: ColorBlue, Body: []Coord{{0, 0}}, Dead: true}
}

func (s *Snake) Head() Coord {
	return s.Body[len(s.Body)-1]
}

func (s *Snake) Eat() {
	s.Len++
}

func (s *Snake) Move() {
	// Priority: Handle target correction if active
	if s.Target != Zero() {
		head := s.Head()
		dx := s.Target.X - head.X
		dy := s.Target.Y - head.Y

		// Calculate half-distance (rounded up, min 1)
		moveX := max(abs(dx)/2, 1) * sign(dx)
		moveY := max(abs(dy)/2, 1) * sign(dy)

		// Clamp to remaining distance (avoid overshooting)
		if abs(moveX) > abs(dx) {
			moveX = dx
		}
		if abs(moveY) > abs(dy) {
			moveY = dy
		}

		newHead := head.Translate(moveX, moveY)
		s.Body = append(s.Body, newHead)

		if len(s.Body) > s.Len {
			s.Body = s.Body[len(s.Body)-s.Len:]
		}

		if newHead == s.Target {
			s.Target = Zero()
		}
		return
	}

	// default movement with speed modifier
	speed := 1
	if s.Speed {
		speed = 2
	}

	for range speed {
		head := s.Head()
		var newHead Coord
		switch s.Dir {
		case UP:
			newHead = head.Translate(0, -1)
		case DOWN:
			newHead = head.Translate(0, 1)
		case LEFT:
			newHead = head.Translate(-1, 0)
		case RIGHT:
			newHead = head.Translate(1, 0)
		}
		s.Body = append(s.Body, newHead)
	}

	if len(s.Body) > s.Len {
		s.Body = s.Body[len(s.Body)-s.Len:]
	}
}

// Helper for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *Snake) ChangeDir(Dir int) bool {
	if ((s.Dir == RIGHT || s.Dir == LEFT) && (Dir == UP || Dir == DOWN)) ||
		((s.Dir == UP || s.Dir == DOWN) && (Dir == LEFT || Dir == RIGHT)) {
		s.Dir = Dir
		return true
	}
	return false
}

func (s *Snake) ChangeSpeed() {
	s.Speed = !s.Speed
	if s.Speed && s.Len <= 2 {
		s.Speed = false
	}
}

func ImportSnake(data []byte) (*Snake, error) {
	var snake Snake
	err := json.Unmarshal(data, &snake)
	if err != nil {
		return nil, err
	}
	return &snake, nil
}

func (s *Snake) Export() ([]byte, error) {

	exportSnake := *s
	size := len(exportSnake.Body)
	last := int(math.Min(1, float64(size)))
	exportSnake.Body = exportSnake.Body[size-last:]

	return json.Marshal(exportSnake)
}
