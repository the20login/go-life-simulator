package dweller

import (
	"time"
	"github.com/the20login/go-life-simulator/world/quadtree"
	"sync"
)

type FoodWorld interface {
	RemoveFood(int32)
	AddFood(quadtree.Point, *Food)
	Size() quadtree.Rectangle
	NextId() int32
}

type FoodProperties struct {
	reproductionRate int
	reproductionRange int
}

var props = FoodProperties{3, 10}

type Food struct {
	Id int32
	world FoodWorld
	die chan int
	Position quadtree.Point
	sync sync.Once
}

func (food *Food) DoAi() {
	for {
		select {
		case <-food.die:
			food.world.RemoveFood(food.Id)
			return
		case <-time.After(30 * time.Second):
			food.reproduce()
		}
	}
}

func (food *Food) Die() {
	food.sync.Do(func(){food.die <- 1})
}

func NewFood(world FoodWorld, position quadtree.Point) *Food {
	return &Food{world.NextId(), world, make(chan int), position, sync.Once{}}
}

func (food *Food) reproduce() {
	position := chooseChildPosition(food.Position, props.reproductionRange, food.world.Size())
	newFood := NewFood(food.world, position)
	food.world.AddFood(position, newFood)
}