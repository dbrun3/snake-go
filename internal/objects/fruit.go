package objects

type Fruit struct {
	Color Color `json:"color"`
}

func CreateFruit(color Color, coord Coord, fruits *map[Coord]Fruit) {
	(*fruits)[coord] = Fruit{Color: color}
}

func DeleteFruit(coord Coord, fruits *map[Coord]Fruit) {
	delete(*fruits, coord)
}
