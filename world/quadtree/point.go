package quadtree

import "math"

type Point struct {
	X, Y float64
}

func (p1 Point) SquareDistance(p2 Point) float64{
	deltaX := p1.X - p2.X
	deltaY := p1.Y - p2.Y
	return deltaX * deltaX + deltaY * deltaY
}

func (p1 Point) Distance(p2 Point) float64 {
	return math.Sqrt(p1.SquareDistance(p2))
}

func (p Point) Delta(dx float64, dy float64) Point {
	return Point{p.X + dx, p.Y + dy}
}

func (p Point) DeltaVector(v Vector) Point{
	return p.Delta(v.dx, v.dy)
}

func (p Point) WithinCircle(point Point, radius float64) bool {
	return p.SquareDistance(point) <= radius * radius
}
