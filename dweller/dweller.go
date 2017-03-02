package dweller

import (
	"github.com/the20login/go-life-simulator/world/quadtree"
	"math/rand"
)

func chooseChildPosition(point quadtree.Point, reproductionRange int, worldSize quadtree.Rectangle) quadtree.Point {
	childVector := quadtree.NewUnitVector(-1 + rand.Float64() * 2).Scale(rand.Float64() * float64(reproductionRange))
	childPoint := point.DeltaVector(childVector)

	if childPoint.X < worldSize.TopLeft().X {
		childPoint = childPoint.Delta(worldSize.TopLeft().X, 0)
	} else if childPoint.X > worldSize.BottomRight().X {
		childPoint = childPoint.Delta(-worldSize.BottomRight().X, 0)
	}

	if childPoint.Y < worldSize.TopLeft().Y {
		childPoint = childPoint.Delta(0, worldSize.TopLeft().Y)
	} else if childPoint.Y > worldSize.BottomRight().Y {
		childPoint = childPoint.Delta(0, -worldSize.BottomRight().Y)
	}
	return childPoint
}