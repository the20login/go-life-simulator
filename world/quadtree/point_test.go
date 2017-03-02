package quadtree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPoint_WithinCircle(t *testing.T) {
	point := Point{50, 57.6}

	assert.False(t, point.WithinCircle(Point{50, 50}, 7.5))
}

var point1 = Point{10, 10}
var point2 = Point{20, 20}
const radius = 15

var result_float float64
func BenchmarkPoint_SquareDistance(b *testing.B) {
	var r float64
	for n := 0; n < b.N; n++ {
		point1.SquareDistance(point2)
	}
	result_float = r
}

var result_bool bool
func BenchmarkPoint_WithinCircle(b *testing.B) {
	var r bool
	for n := 0; n < b.N; n++ {
		r = point1.WithinCircle(point2, radius)
	}
	result_bool = r
}