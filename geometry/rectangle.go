package geometry

import "fmt"

// =================================================================
// Rectangle
// =================================================================

type Rectangle struct {
	X, Y, Width, Height float32
}

// AABB returns the axis-aligned bounding box of the rectangle.
func (r Rectangle) AABB() Rectangle {
	return r
}

// Contains checks if the rectangle contains a given point.
func (r Rectangle) Contains(point Point) bool {
	return point.X >= r.X && point.X <= r.X+r.Width && point.Y >= r.Y && point.Y <= r.Y+r.Height
}

// Intersects checks if the rectangle intersects with another rectangle.
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.X < other.X+other.Width && r.X+r.Width > other.X && r.Y < other.Y+other.Height && r.Y+r.Height > other.Y
}

// IntersectsShape checks if the rectangle intersects with a given shape's axis-aligned bounding box.
func (r Rectangle) IntersectsShape(shape Shape) bool {
	// Check if the rectangle intersects with the axis-aligned bounding box of the shape
	return r.Intersects(shape.AABB())
}

// String returns a string representation of the rectangle.
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(X: %.2f, Y: %.2f, Width: %.2f, Height: %.2f)", r.X, r.Y, r.Width, r.Height)
}
