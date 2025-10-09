package quadtree

import (
	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-core/hashset"
)

func (qt *QuadTree[T]) GetQuadTreeNode(x, y, width, height float64, leafSize, depth int) *QuadTree[T] {
	for node := range qt.nodePool {
		qt.nodePool.Remove(node)
		node.bounds.X = x
		node.bounds.Y = y
		node.bounds.Width = width
		node.bounds.Height = height
		node.leafSize = leafSize
		node.depth = depth
		return node
	}
	return New[T](geom.NewRect64(x, y, width, height), leafSize, depth)
}

func (qt *QuadTree[T]) PutQuadTreeNode(node *QuadTree[T]) {
	node.Clear()
	qt.nodePool.AddDistinct(node)
}

type QuadTree[T geom.Bounded] struct {
	bounds   geom.Rect64
	objects  hashset.Set[T]
	nodes    [4]*QuadTree[T]
	nodePool hashset.Set[*QuadTree[T]]

	leafSize int
	depth    int
}

func New[T geom.Bounded](bounds geom.Rect64, leafSize, depth int) *QuadTree[T] {
	return &QuadTree[T]{
		bounds:   bounds,
		objects:  hashset.New[T](),
		nodes:    [4]*QuadTree[T]{},
		nodePool: hashset.New[*QuadTree[T]](),
		leafSize: leafSize,
		depth:    depth,
	}
}

func (qt *QuadTree[T]) Count() int {
	count := len(qt.objects)
	if qt.isBranch() {
		for _, node := range qt.nodes {
			count += node.Count()
		}
	}
	return count
}

func (qt *QuadTree[T]) Resize(bounds geom.Rect64) {
	if qt.bounds == bounds {
		return
	}

	allObjects := hashset.New[T]()

	qt.internalQuery(qt.bounds, allObjects)
	qt.Clear()

	qt.bounds = bounds
	for o := range allObjects {
		qt.Insert(o)
	}
}

func (qt *QuadTree[T]) Min() (float64, float64) {
	return qt.bounds.Min()
}

func (qt *QuadTree[T]) Max() (float64, float64) {
	return qt.bounds.Max()
}

func (qt *QuadTree[T]) Insert(obj T) bool {
	bounds := obj.Bounds()
	if !bounds.Intersects(qt.bounds) {
		return false
	}

	if qt.isBranch() {
		for _, node := range qt.nodes {
			node.Insert(obj)
		}
		return true
	}

	// If this leaf is full, then allow it to overflow if we are at max depth.
	if len(qt.objects) < qt.leafSize || qt.depth <= 0 {
		qt.objects.AddDistinct(obj)
		return true
	}

	qt.subdivide()
	for o := range qt.objects {
		for _, node := range qt.nodes {
			node.Insert(o)
		}
	}
	qt.objects.Clear()

	added := false
	for _, node := range qt.nodes {
		if node.Insert(obj) {
			added = true
		}
	}
	return added
}

func (qt *QuadTree[T]) Remove(obj T) bool {
	removed := false

	if qt.objects.Contains(obj) {
		qt.objects.Remove(obj)
		removed = true
	}

	if qt.isBranch() {
		for _, node := range qt.nodes {
			if node.Remove(obj) {
				removed = true
			}
		}
		qt.tryCollapse()
	}

	return removed
}

func (qt *QuadTree[T]) Update(obj T) {
	qt.Remove(obj)
	qt.Insert(obj)
}

func (qt *QuadTree[T]) Query(region geom.Rect64) hashset.Set[T] {
	results := hashset.New[T]()
	qt.internalQuery(region, results)
	return results
}

func (qt *QuadTree[T]) QueryNodes(region geom.Rect64) hashset.Set[*QuadTree[T]] {
	results := hashset.New[*QuadTree[T]]()
	qt.internalQueryNodes(region, results)
	return results
}

func (qt *QuadTree[T]) Partitions(area geom.Rect64) hashset.Set[geom.Rect64] {
	results := hashset.New[geom.Rect64]()
	qt.internalPartitions(results)
	return results
}

// Clear removes all objects and child nodes from the quadtree, resetting it to an empty state.
//
// Clearing the quadtree only removes references of the objects stored within it; it does not
// delete the objects themselves.
func (qt *QuadTree[T]) Clear() {
	qt.objects.Clear()
	if qt.isBranch() {
		for i, node := range qt.nodes {
			qt.PutQuadTreeNode(node)
			qt.nodes[i] = nil
		}
	}
}

func (qt *QuadTree[T]) isBranch() bool {
	// Note: As long as leafs are managed correctly, checking the first node is sufficient
	return qt.nodes[0] != nil
}

func (qt *QuadTree[T]) subdivide() {
	halfWidth := qt.bounds.Width / 2
	halfHeight := qt.bounds.Height / 2

	x := qt.bounds.X
	y := qt.bounds.Y

	qt.nodes[0] = qt.GetQuadTreeNode(x, y, halfWidth, halfHeight, qt.leafSize, qt.depth-1)                      // NW
	qt.nodes[1] = qt.GetQuadTreeNode(x+halfWidth, y, halfWidth, halfHeight, qt.leafSize, qt.depth-1)            // NE
	qt.nodes[2] = qt.GetQuadTreeNode(x, y+halfHeight, halfWidth, halfHeight, qt.leafSize, qt.depth-1)           // SW
	qt.nodes[3] = qt.GetQuadTreeNode(x+halfWidth, y+halfHeight, halfWidth, halfHeight, qt.leafSize, qt.depth-1) // SE
}

func (qt *QuadTree[T]) tryCollapse() {
	if !qt.isBranch() {
		return
	}

	allObjects := hashset.New[T]()
	qt.internalQuery(qt.bounds, allObjects)

	if len(allObjects) <= qt.leafSize {
		qt.objects = allObjects
		for i, node := range qt.nodes {
			qt.PutQuadTreeNode(node)
			qt.nodes[i] = nil
		}
	}
}

func (qt *QuadTree[T]) internalQuery(region geom.Rect64, results hashset.Set[T]) {
	if !region.Intersects(qt.bounds) {
		return
	}

	for o := range qt.objects {
		results.AddDistinct(o)
	}

	if qt.isBranch() {
		for _, node := range qt.nodes {
			node.internalQuery(region, results)
		}
	}
}

func (qt *QuadTree[T]) internalQueryNodes(region geom.Rect64, results hashset.Set[*QuadTree[T]]) {
	if !region.Intersects(qt.bounds) {
		return
	}

	results.AddDistinct(qt)

	if qt.isBranch() {
		for _, node := range qt.nodes {
			node.internalQueryNodes(region, results)
		}
	}
}

func (qt *QuadTree[T]) internalPartitions(results hashset.Set[geom.Rect64]) {
	results.AddDistinct(qt.bounds)
	if qt.isBranch() {
		for _, node := range qt.nodes {
			node.internalPartitions(results)
		}
	}
}
