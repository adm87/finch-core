package geometry

// Shape interface defines the methods that any geometric shape must implement.
type Shape interface {
	// AABB returns the axis-aligned bounding box of the shape.
	AABB() Rectangle

	// Contains checks if the shape contains a given point.
	Contains(point Point) bool

	// Intersects checks if the shape intersects with another shape.
	Intersects(other Shape) bool
}
