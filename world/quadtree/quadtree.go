package quadtree

import "fmt"

type QuadTree struct {
	count int
	root  *Node
}

func NewQuadTree(rectangle Rectangle) *QuadTree {
	return &QuadTree{0, NewNode(rectangle, nil)}
}

func (qt *QuadTree) Put(point Point, value int32) {
	if !qt.root.rectangle.Contains(point) {
		panic(fmt.Sprint("Out of bounds: ", point))
	}
	if insert(qt.root, point, value) {
		qt.count++;
	}
}

/**
     * Gets the value of the point
     * <br>
     * Note: this method require absolute precision({@link Point#equals Point.equals}), use {@link #searchWithin} to search with custom precision
     *
     * @param point coordinates
     * @return Optional value at point
     */
func (qt *QuadTree) Get(point Point) (int32, bool) {
	node, ok := find(qt.root, point)
	if ok {
		return node.value, true
	} else {
		return -1, false
	}
}

/**
 * Removes element from tree. Tree will rebalance itself after removing.
 * <br>
 * Note: this method require absolute precision({@link Point#equals Point.equals})
 *
 * @param point coordinates
 * @return Optional removed element, empty if there is no element at provided coordinates
 */
func (qt *QuadTree) Remove(point Point) (int32, bool) {
	node, ok := find(qt.root, point)
	if ok {
		value := node.value
		node.clear()
		balance(node)
		qt.count--
		return value, true
	} else {
		return -1, false
	}
}

/**
 * Checks if there is an element at this point
 *
 * @param point coordinates
 * @return True if tree Contains element at this point, otherwise false.
 */
func (qt *QuadTree) Contains(point Point) bool {
	_, ok := qt.Get(point)
	return ok
}

/**
 * @return Whether the tree is empty.
 */
func (qt *QuadTree) IsEmpty() bool {
	//return qt.root.nodeType == empty;
	return qt.count == 0
}

/**
 * @return The number of elements in the tree.
 */
func (qt *QuadTree) Count() int {
	return qt.count
}

/**
 * Removes all elements from the tree.
 */
func (qt *QuadTree) Clear() {
	qt.root.clear()
	qt.count = 0
}

/**
 * Returns all elements of the tree
 * @return Stream of tree elements
 */
func (qt *QuadTree) Values() []int32 {
	var values []int32
	traverse(qt.root, func(node *Node) { values = append(values, node.value) })
	return values
}

/**
 * Returns elements within rectangle(inclusive)
 * @return Stream of tree elements
 */
func (qt *QuadTree) SearchWithinRectangle(rectangle Rectangle) []int32 {
	var values []int32
	navigateRectangle(qt.root, rectangle, func(value int32) { values = append(values, value) })
	return values
}

/**
 * Returns elements within circle(inclusive)
 * @return Stream of tree elements
 */
func (qt *QuadTree) SearchWithinCircle(point Point, radius float64) []int32 {
	values := make([]int32, 0, 32)
	navigateCircle(qt.root, point, radius, func(value int32) { values = append(values, value) })
	return values
}

func navigateCircle(node *Node, point Point, radius float64, consumer func(int32)) {
	switch (node.nodeType) {
	case leaf:
		if node.point.WithinCircle(point, radius) {
			consumer(node.value)
		}
	case pointer:
		if node.neNode.rectangle.IsIntersectCircle(point, radius) {
			navigateCircle(node.neNode, point, radius, consumer)
		}
		if node.seNode.rectangle.IsIntersectCircle(point, radius) {
			navigateCircle(node.seNode, point, radius, consumer)
		}
		if node.swNode.rectangle.IsIntersectCircle(point, radius) {
			navigateCircle(node.swNode, point, radius, consumer)
		}
		if node.nwNode.rectangle.IsIntersectCircle(point, radius) {
			navigateCircle(node.nwNode, point, radius, consumer)
		}
	}
}

func navigateRectangle(node *Node, rectangle Rectangle, consumer func(int32)) {
	switch node.nodeType {
	case leaf:
		if rectangle.Contains(node.point) {
			consumer(node.value)
		}
	case pointer:
		if node.neNode.rectangle.IsIntersectRectangle(rectangle) {
			navigateRectangle(node.neNode, rectangle, consumer)
		}
		if node.seNode.rectangle.IsIntersectRectangle(rectangle) {
			navigateRectangle(node.seNode, rectangle, consumer)
		}
		if node.swNode.rectangle.IsIntersectRectangle(rectangle) {
			navigateRectangle(node.swNode, rectangle, consumer)
		}
		if node.nwNode.rectangle.IsIntersectRectangle(rectangle) {
			navigateRectangle(node.nwNode, rectangle, consumer)
		}
	}
}

/**
 * Clones the quad-tree and returns the new instance.
 * @return {QuadTree} A clone of the tree.
 */
func (qt *QuadTree) clone() *QuadTree {
	clone := NewQuadTree(qt.root.rectangle);
	// This is inefficient as the clone needs to recalculate the structure of the
	// tree, even though we know it already.  But this is easier and can be
	// optimized when/if needed.
	traverse(qt.root, func(node *Node) { clone.Put(node.point, node.value) })

	return clone
}

func traverse(node *Node, consumer func(*Node)) {
	switch (node.nodeType) {
	case leaf:
		consumer(node)
	case pointer:
		for _, child := range node.getChildNodes() {
			traverse(child, consumer)
		}
	}
}

func find(node *Node, point Point) (*Node, bool) {
	switch (node.nodeType) {
	case empty:
		return nil, false;
	case leaf:
		if node.point == point {
			return node, true
		} else {
			return nil, false
		}
	case pointer:
		return find(node.getQuadrantNode(point), point)
	default:
		panic("Invalid node type")
	}
}

func insert(parent *Node, point Point, value int32) bool {
	switch (parent.nodeType) {
	case empty:
		setPointForNode(parent, point, value)
		return true
	case leaf:
		if parent.point == point {
			setPointForNode(parent, point, value)
			return false
		} else {
			split(parent)
			return insert(parent.getQuadrantNode(point), point, value)
		}
	case pointer:
		return insert(parent.getQuadrantNode(point), point, value)
	default:
		panic("Invalid nodeType in parent")
	}
}

func split(node *Node) {
	oldPoint := node.point
	oldValue := node.value

	node.clearAndSplit()

	insert(node, oldPoint, oldValue);
}

func balance(node *Node) {
	switch node.nodeType {
	case empty:
		fallthrough
	case leaf:
		if node.parent != nil {
			balance(node.parent)
		}
	case pointer:
		var nonEmptyNodes []*Node
		for _, child := range node.getChildNodes() {
			if child.nodeType != empty {
				nonEmptyNodes = append(nonEmptyNodes, child)
			}
		}

		switch len(nonEmptyNodes) {
		case 0:
			node.clear()
		case 1:
			firstLeaf := nonEmptyNodes[0];
			if (firstLeaf.nodeType == pointer) {
				// Only child was a pointer, therefore we can't rebalance.
				break
			} else {
				// Only child was a leaf: so update node's point and make it a leaf.
				node.setPoint(firstLeaf.point, firstLeaf.value)
			}
		}

		// Try and balance the parent as well.
		if node.parent != nil {
			balance(node.parent)
		}
	}
}

func setPointForNode(node *Node, point Point, value int32) {
	if (node.nodeType == pointer) {
		panic("Can not put point for node of type POINTER")
	}
	node.setPoint(point, value)
}
