package quadtree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math/rand"
	/*"sync"
	"net/http"
	 _ "net/http/pprof"*/
)

func create() *QuadTree {
	quadTree := NewQuadTree(Rectangle{Point{0, 0}, Point{100, 100}})
	quadTree.Put(Point{0, 0}, 0)
	quadTree.Put(Point{0, 2}, 1)
	quadTree.Put(Point{2, 0}, 2)
	quadTree.Put(Point{2, 2}, 3)
	quadTree.Put(Point{100, 100}, 4)
	quadTree.Put(Point{1, 90}, 5)
	quadTree.Put(Point{50, 50}, 6)
	quadTree.Put(Point{8, 54}, 7)
	quadTree.Put(Point{55, 45}, 8)
	quadTree.Put(Point{45, 55}, 9)
	quadTree.Put(Point{50, 57.6}, 10)
	return quadTree
}

func TestQuadTree_SearchWithinRectangle(t *testing.T) {
	qt := create()
	values := qt.SearchWithinRectangle(Rectangle{Point{0, 0}, Point{10, 10}})

	assert.Equal(t, 4, len(values))
	assert.Contains(t, values, int32(0))
	assert.Contains(t, values, int32(1))
	assert.Contains(t, values, int32(2))
	assert.Contains(t, values, int32(3))

}

func TestQuadTree_SearchWithinCircle(t *testing.T) {
	qt := create()
	values := qt.SearchWithinCircle(Point{50, 50}, 7.5)

	assert.Equal(t, 3, len(values))
	assert.Contains(t, values, int32(6))
	assert.Contains(t, values, int32(8))
	assert.Contains(t, values, int32(9))
}

func TestQuadTree_SearchWithinRectangle_NoResult(t *testing.T) {
	qt := create()
	values := qt.SearchWithinRectangle(Rectangle{Point{0.5, 0.5}, Point{1.5, 1.5}})

	assert.Equal(t, 0, len(values))
}

func TestQuadTree_SearchWithinCircle_NoResult(t *testing.T) {
	qt := create()
	values := qt.SearchWithinCircle(Point{1, 1}, 0.7)

	assert.Equal(t, 0, len(values))
}




var random = rand.New(rand.NewSource(777))
const worldSize = 10000
const points = 100000;
var center = Point{worldSize / 2, worldSize / 2}
var searchRadius = 1000


func generate() *QuadTree {
	qt := NewQuadTree(Rectangle{Point{0, 0}, Point{worldSize, worldSize}})

	for i := int32(0); i< points; i++ {
		qt.Put(Point{random.Float64() * worldSize, random.Float64() * worldSize}, i)
	}
	return qt
}

var searchResult int
func BenchmarkQuadTree_SearchWithinRectangle(b *testing.B) {
	qt := generate()
	rectangle := Rectangle{Point{center.X - float64(searchRadius), center.Y - float64(searchRadius)}, Point{center.X + float64(searchRadius), center.Y + float64(searchRadius)}}

	b.ResetTimer()
	var r int
	for n := 0; n < b.N; n++ {
		r = len(qt.SearchWithinRectangle(rectangle))
	}
	searchResult = r
}

/*var once= sync.Once{}
func TestQuadTree_SearchWithinCircle2(t *testing.T) {
	once.Do(func() {go http.ListenAndServe("localhost:6060", nil)})
	qt := generate()
	var r int
	for n := 0; n <1e15; n++ {
		r = len(qt.SearchWithinCircle(center, float64(searchRadius)))
	}
	searchResult = r
}*/

func BenchmarkQuadTree_SearchWithinCircle(b *testing.B) {
	qt := generate()
	b.ResetTimer()
	var r int
	for n := 0; n < b.N; n++ {
		r = len(qt.SearchWithinCircle(center, float64(searchRadius)))
	}
	searchResult = r
}