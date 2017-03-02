package world

import (
	"github.com/the20login/go-life-simulator/world/quadtree"
	"github.com/the20login/go-life-simulator/dweller"
	"sync/atomic"
	"sync"
	"fmt"
)

type World struct {
	Rectangle     quadtree.Rectangle
	foodPositions *quadtree.QuadTree
	antPositions  *quadtree.QuadTree
	foodMap       map[int32]*dweller.Food
	antsMap       map[int32]*dweller.Ant
	idGenerator   int32
	foodLock      sync.RWMutex
	antsLock      sync.RWMutex
}

func NewWorld(rectangle quadtree.Rectangle) *World {
	return &World{rectangle, quadtree.NewQuadTree(rectangle), quadtree.NewQuadTree(rectangle), make(map[int32]*dweller.Food), make(map[int32]*dweller.Ant), 1, sync.RWMutex{}, sync.RWMutex{}}}

func (w *World) Size() quadtree.Rectangle {
	return w.Rectangle
}

func (w *World) NextId() int32 {
	return atomic.AddInt32(&w.idGenerator, 1)
}

func (w *World) AddFood(point quadtree.Point, food *dweller.Food) {
	w.foodLock.Lock()
	defer w.foodLock.Unlock()

	food.Position = point
	w.foodMap[food.Id] = food
	w.foodPositions.Put(point, food.Id)

	fmt.Println("new food added with id", food.Id)

	go food.DoAi()
}

func (w *World) RemoveFood(id int32) {
	w.foodLock.Lock()
	defer w.foodLock.Unlock()

	position := w.foodMap[id].Position
	w.foodPositions.Remove(position)
	delete(w.foodMap, id)
}

func (w *World) AddAnt(point quadtree.Point, ant *dweller.Ant) {
	w.antsLock.Lock()
	defer w.antsLock.Unlock()

	ant.Position = point
	w.antsMap[ant.Id] = ant
	w.antPositions.Put(point, ant.Id)

	fmt.Println("new ant added with id", ant.Id)

	go ant.DoAi()
}

func (w *World) RemoveAnt(id int32) {
	w.antsLock.Lock()
	defer w.antsLock.Unlock()

	position := w.antsMap[id].Position
	w.antPositions.Remove(position)
	delete(w.antsMap, id)
}

func (w *World) MoveAnt(id int32, target quadtree.Point) {
	w.antsLock.Lock()
	defer w.antsLock.Unlock()

	ant := w.antsMap[id]
	w.antPositions.Remove(ant.Position)
	w.antPositions.Put(target, id)
	ant.Position = target
}

func (w *World) SearchFoodInRange(point quadtree.Point, radius float64) []*dweller.Food {
	w.foodLock.RLock()
	defer w.foodLock.RUnlock()

	var array []*dweller.Food
	for _, id := range w.foodPositions.SearchWithinCircle(point, radius) {
		array = append(array, w.foodMap[id])
	}
	return array
}

func (w *World) FoodPositions() []quadtree.Point {
	w.foodLock.RLock()
	defer w.foodLock.RUnlock()

	var array []quadtree.Point
	for _, food := range w.foodMap {
		array = append(array, food.Position)
	}
	return array
}

func (w *World) AntsPositions() []quadtree.Point {
	w.antsLock.RLock()
	defer w.antsLock.RUnlock()

	var array []quadtree.Point
	for _, ant := range w.antsMap {
		array = append(array, ant.Position)
	}
	return array
}