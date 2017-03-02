package quadtree

const (
	empty = "empty"
	leaf = "leaf"
	pointer = "pointer"
)

type Node struct {
	rectangle Rectangle
	parent *Node
	point Point
	nodeType string
	value int32
	nwNode *Node
	neNode *Node
	swNode *Node
	seNode *Node
}

func NewNode(rectangle Rectangle, parent *Node) *Node {
	return &Node{rectangle:rectangle, parent:parent, nodeType:empty, value:-1}
}

func (node *Node) setPoint(point Point, payload int32) {
	node.clear()
	node.nodeType = leaf
	node.point = point
	node.value = payload
}

func (node *Node) clear() {
	node.nodeType = empty
	node.nwNode = nil
	node.neNode = nil
	node.swNode = nil
	node.seNode = nil
	node.value = -1
}

func (node *Node) clearAndSplit() {
	node.clear()
	node.nodeType = pointer

	center := node.rectangle.Center()
	halfWidth := node.rectangle.Width() / 2
	halfHeight := node.rectangle.Height() / 2

	node.nwNode = NewNode(Rectangle{node.rectangle.TopLeft(), center}, node)
	node.neNode = NewNode(Rectangle{center.Delta(0, -halfHeight), center.Delta(halfWidth, 0)}, node)
	node.swNode = NewNode(Rectangle{center.Delta(-halfWidth, 0), center.Delta(0, halfHeight)}, node)
	node.seNode = NewNode(Rectangle{center, node.rectangle.BottomRight()}, node)
}

func (node *Node) getQuadrantNode(point Point) *Node {
	middle := node.rectangle.Center()
	if point.X < middle.X {
		if point.Y < middle.Y {
			return node.nwNode
		} else {
			return node.swNode
		}
	} else {
		if point.Y < middle.Y {
			return node.neNode
		} else {
			return node.seNode
		}
	}
}

func (node *Node) getChildNodes() []*Node {
	return []*Node{node.neNode, node.seNode, node.swNode, node.nwNode}
}
