package quadtree

type Quadrant int

const (
	NW Quadrant = iota
	NE Quadrant = iota
	SW Quadrant = iota
	SE Quadrant = iota
)
type Rectangle struct {
	p1, p2 Point
}

func NewRectangle(topLeft Point, bottomRight Point) Rectangle {
	return Rectangle{topLeft, bottomRight}
}

func (r Rectangle) TopLeft() Point {
	return r.p1
}

func (r Rectangle) TopRight() Point {
	return Point{r.p2.X, r.p1.Y}
}

func (r Rectangle) BottomLeft() Point {
return Point{r.p1.X, r.p2.Y}
}

func (r Rectangle) BottomRight() Point {
	return r.p2
}

func (r Rectangle) Center() Point{
return Point{r.p1.X + r.Width() / 2, r.p1.Y + r.Height() / 2}
}

func (r Rectangle) Width() float64{
	return r.p2.X - r.p1.X
}

func (r Rectangle) Height() float64{
	return r.p2.Y - r.p1.Y
}

func (r Rectangle) Contains(point Point) bool{
	return r.p1.X <= point.X && r.p2.X >= point.X && r.p1.Y <= point.Y && r.p2.Y >= point.Y
}

func (r Rectangle) IsIntersectRectangle(rect Rectangle) bool {
	return r.p1.X <= rect.p2.X && r.p2.X >= rect.p1.X && r.p1.Y <= rect.p2.Y && r.p2.Y >= rect.p1.Y
}

func (r Rectangle) IsIntersectCircle(point Point, radius float64) bool {
	var deltaX, deltaY float64
	if point.X < r.p1.X {
		deltaX = point.X - r.p1.X
	} else if point.X > r.p2.X {
		deltaX = point.X - r.p2.X
	} else {
		deltaX = 0
	}

	if point.Y < r.p1.Y {
		deltaY = point.Y - r.p1.Y
	} else if point.Y > r.p2.Y {
		deltaY = point.Y - r.p2.Y
	} else {
		deltaY = 0
	}

	return (deltaX * deltaX + deltaY * deltaY) < (radius * radius)
}

func (r Rectangle) pointInQuadrant(point Point) Quadrant {
	center := r.Center()
	if point.X < center.X {
		if point.Y < center.Y {return NW} else {return SW}
	} else {
		if point.Y < center.Y {return NE} else {return SE}
	}
}

func (r Rectangle) fitInRectangle(rectangle Rectangle) bool {
	return r.p1.X >= rectangle.p1.X && r.p2.X <= rectangle.p2.X && r.p1.Y >= rectangle.p1.Y && r.p2.Y <= rectangle.p2.Y;
}

func (r Rectangle) fitInCircle(point Point, radius float64) bool {
	squareRadius := radius * radius
	return r.TopLeft().SquareDistance(point) <= squareRadius &&
		r.TopRight().SquareDistance(point) <= squareRadius &&
		r.BottomLeft().SquareDistance(point) <= squareRadius &&
		r.BottomRight().SquareDistance(point) <= squareRadius;
}