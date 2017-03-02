package quadtree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
)

func TestVector_NewUnitVector(t *testing.T) {
	v := NewUnitVector(0)
	assert.True(t, math.Abs(float64(1) - v.dx) < 0.000001)
	assert.True(t, math.Abs(float64(0) - v.dy) < 0.000001)

	v = NewUnitVector(math.Pi /2)
	assert.True(t, math.Abs(float64(0) - v.dx) < 0.000001)
	assert.True(t, math.Abs(float64(1) - v.dy) < 0.000001)

	v = NewUnitVector(math.Pi)
	assert.True(t, math.Abs(float64(-1) - v.dx) < 0.000001)
	assert.True(t, math.Abs(float64(0) - v.dy) < 0.000001)

	v = NewUnitVector(math.Pi *3/2)
	assert.True(t, math.Abs(float64(0) - v.dx) < 0.000001)
	assert.True(t, math.Abs(float64(-1) - v.dy) < 0.000001)

	v = NewUnitVector(math.Pi*2)
	assert.True(t, math.Abs(float64(1) - v.dx) < 0.000001)
	assert.True(t, math.Abs(float64(0) - v.dy) < 0.000001)

}