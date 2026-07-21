package mover

type Point struct {
	X float64
	Y float64
}

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func TargetPosition(position Point, distance int, displays []Rect) (Point, bool) {
	if distance <= 0 {
		return Point{}, false
	}
	offset := float64(distance)
	candidates := []Point{
		{X: position.X + offset, Y: position.Y},
		{X: position.X - offset, Y: position.Y},
		{X: position.X, Y: position.Y + offset},
		{X: position.X, Y: position.Y - offset},
	}
	for _, candidate := range candidates {
		if ContainsAny(displays, candidate) {
			return candidate, true
		}
	}
	return Point{}, false
}

func ContainsAny(displays []Rect, point Point) bool {
	for _, display := range displays {
		if display.Contains(point) {
			return true
		}
	}
	return false
}

func (r Rect) Contains(point Point) bool {
	return point.X >= r.X &&
		point.X < r.X+r.Width &&
		point.Y >= r.Y &&
		point.Y < r.Y+r.Height
}
