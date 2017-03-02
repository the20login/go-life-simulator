package dweller

import (
	"time"
	"github.com/the20login/go-life-simulator/world/quadtree"
	"github.com/bradfitz/slice"
	"math/rand"
	"fmt"
	"math"
)

type AntWorld interface {
	RemoveAnt(int32)
	AddAnt(quadtree.Point, *Ant)
	Size() quadtree.Rectangle
	NextId() int32
	SearchFoodInRange(point quadtree.Point, radius float64) []*Food
	MoveAnt(int32, quadtree.Point)
}

type AntProperties struct {
	visibilityRange int
	actionRange int
	speed int
}

var propsAnts = AntProperties{20, 2, 5}

type Ant struct {
	Id           int32
	world        AntWorld
	die          chan int
	Position     quadtree.Point
	targetFood   *Food
	randomVector quadtree.Vector
	canBreed bool
	breedCooldown <-chan time.Time
}

func (ant *Ant) DoAi() {
	ant.breedCooldown = time.After(10*time.Second)
	for {
		select {
		case <-ant.die:
			ant.world.RemoveAnt(ant.Id)
			return
		case <- ant.breedCooldown:
			ant.canBreed = true
		case <-time.After(time.Second):
			if ant.canBreed {
				ant.reproduce()
			}
			target := ant.selectTarget()
			ant.world.MoveAnt(ant.Id, target)
		}
	}
}

func (ant *Ant) Die() {
	ant.die <- 1
}

func NewAnt(world AntWorld, position quadtree.Point) *Ant {
	return &Ant{world.NextId(), world, make(chan int), position, nil, quadtree.NewZeroVector(), false, nil}
}

func (ant *Ant) reproduce() {
	position := chooseChildPosition(ant.Position, props.reproductionRange, ant.world.Size())
	newFood := NewAnt(ant.world, position)
	ant.world.AddAnt(position, newFood)
	ant.canBreed = false
	ant.breedCooldown = time.After(20*time.Second)
}

func (ant *Ant) selectTarget() quadtree.Point {
	if ant.targetFood == nil {
		ant.targetFood = ant.chooseFood()
	}

	if ant.targetFood == nil {
		if ant.randomVector.IsZeroVector() {
			ant.randomVector = quadtree.NewUnitVector(rand.Float64() * 2 * math.Pi).Scale(rand.Float64() * 5)
		}
		return ant.Position.DeltaVector(ant.randomVector)
	} else {
		ant.randomVector = quadtree.NewZeroVector()
		if ant.targetFood.Position.SquareDistance(ant.Position) <= float64(propsAnts.actionRange * propsAnts.actionRange) {
			ant.targetFood.Die()
			fmt.Println("food", ant.targetFood.Id, "eaten by", ant.Id)
			ant.targetFood = nil
			return ant.Position
		}
		return ant.calculateMove(ant.targetFood.Position)
	}
}

func (ant *Ant) calculateMove(target quadtree.Point) quadtree.Point {
	targetVector := quadtree.NewVector(ant.Position, target)
	var speedVector quadtree.Vector
	if targetVector.SquareLength() <= float64(propsAnts.speed*propsAnts.speed) {
		length := targetVector.Length()
		speedVector = targetVector.Scale((length - float64(propsAnts.actionRange)/2) / length)
	} else {
		speedVector = targetVector.Scale(float64(propsAnts.speed) / targetVector.Length())
	}

	return ant.Position.DeltaVector(speedVector)
}

func (ant *Ant) chooseFood() *Food {
	food := ant.world.SearchFoodInRange(ant.Position, float64(propsAnts.visibilityRange))
	if len(food) == 0 {
		return nil
	}
	slice.Sort(food, func(i, j int) bool {return food[i].Position.SquareDistance(ant.Position) < food[j].Position.SquareDistance(ant.Position)})

	return food[0]
}
