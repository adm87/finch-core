package geometry

// Intersects check if the two shapes intersect.
//
// It first checks if the axis-aligned bounding boxes (AABBs) of the shapes intersect before checking for intersection between the shapes themselves.
func Intersects(shapeA, shapeB Shape) bool {
	aabbA := shapeA.AABB()
	aabbB := shapeB.AABB()

	if !aabbA.Intersects(aabbB) {
		return false
	}

	return shapeA.Intersects(shapeB)
}
