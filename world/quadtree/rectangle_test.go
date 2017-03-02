package quadtree

import "testing"



var result bool
func BenchmarkRectangle_IsIntersectCircle_Corner(b *testing.B) {
	rect := Rectangle{Point{10, 10}, Point{100, 100}}
	center := Point{0, 0}
	radius := float64(20)
	b.ResetTimer()

	var r bool
	for i:=0;i<b.N;i++ {
		r = rect.IsIntersectCircle(center, radius)
	}
	result = r
}

func BenchmarkRectangle_IsIntersectCircle_None(b *testing.B) {
	rect := Rectangle{Point{10, 10}, Point{100, 100}}
	center := Point{0, 0}
	radius := float64(10)
	b.ResetTimer()

	var r bool
	for i:=0;i<b.N;i++ {
		r = rect.IsIntersectCircle(center, radius)
	}
	result = r
}

func BenchmarkRectangle_IsIntersectCircle_Inside(b *testing.B) {
	rect := Rectangle{Point{10, 10}, Point{100, 100}}
	center := Point{20, 20}
	radius := float64(10)
	b.ResetTimer()

	var r bool
	for i:=0;i<b.N;i++ {
		r = rect.IsIntersectCircle(center, radius)
	}
	result = r
}