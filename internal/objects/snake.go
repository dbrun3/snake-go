package objects

import (
	"encoding/json"
	"math"
)

const MARSHAL_BODY_SIZE = 5

const (
	UP int = iota
	DOWN
	LEFT
	RIGHT
)

type Snake struct {
	Dir   int     `json:"Dir"`
	Len   int     `json:"Len"`
	Color Color   `json:"Color"`
	Name  string  `json:"Name"`
	Id    string  `json:"Id"`
	Body  []Coord `json:"Body"`
	Dead  bool    `json:"Dead"`
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
	var newHead Coord
	head := s.Head()

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

	// append new head
	s.Body = append(s.Body, newHead)

	if len(s.Body) > s.Len {
		s.Body = s.Body[len(s.Body)-s.Len:]
	}
}

func (s *Snake) ChangeDir(Dir int) bool {
	if ((s.Dir == RIGHT || s.Dir == LEFT) && (Dir == UP || Dir == DOWN)) ||
		((s.Dir == UP || s.Dir == DOWN) && (Dir == LEFT || Dir == RIGHT)) {
		s.Dir = Dir
		return true
	}
	return false
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
	last := int(math.Min(MARSHAL_BODY_SIZE, float64(size)))
	exportSnake.Body = exportSnake.Body[size-last:]

	return json.Marshal(exportSnake)
}
