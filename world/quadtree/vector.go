package quadtree

import (
	"math"
)

type Vector struct {
	dx, dy float64
}

func NewVector(from, to Point) Vector {
	return Vector{to.X - from.X, to.Y - from.Y}
}

func NewUnitVector(angle float64) Vector {
	return Vector{math.Cos(angle), math.Sin(angle)}
}

func NewZeroVector() Vector {
	return Vector{0, 0}
}

func (v Vector) Length() float64{
	return math.Sqrt(v.dx*v.dx + v.dy*v.dy)
}

func (v Vector) SquareLength() float64 {
	return v.dx*v.dx + v.dy*v.dy
}

func (v Vector) Scale(scaleFactor float64) Vector {
	return Vector{v.dx * scaleFactor, v.dy * scaleFactor}
}

func (v Vector) Angle() float64{
	return math.Atan2(v.dy, v.dx);
}

func (v Vector) Rotate(newAngle float64) Vector {
	length := v.Length()
	return Vector{length * math.Cos(newAngle), length * math.Sin(newAngle)}
}

func (v Vector) IsZeroVector() bool {
	return v.dx== 0 && v.dy == 0;
}

func (v Vector) Plus(other Vector) Vector {
	return Vector{v.dx + other.dx, v.dy + other.dy};
}