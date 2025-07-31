package geometry

import (
	"fmt"
	"math"
)

// Point represents a 2D point in space.
type Point struct {
	X, Y float32
}

// Add adds another point to this point.
func (p *Point) Add(other Point) {
	p.X += other.X
	p.Y += other.Y
}

// Subtract subtracts another point from this point.
func (p *Point) Subtract(other Point) {
	p.X -= other.X
	p.Y -= other.Y
}

// Scale scales the point by a scalar value.
func (p *Point) Scale(scalar float32) {
	p.X *= scalar
	p.Y *= scalar
}

// DistanceTo calculates the distance from this point to another point.
func (p Point) DistanceTo(other Point) float32 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Length calculates the length of the point vector from the origin (0, 0).
func (p Point) Length() float32 {
	return float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
}

// Normalize return a normalized version of the point vector.
func (p Point) Normalize() Point {
	length := p.Length()
	if length == 0 {
		return Point{0, 0}
	}
	return Point{p.X / length, p.Y / length}
}

func (p Point) String() string {
	return fmt.Sprintf("Point(X: %.2f, Y: %.2f)", p.X, p.Y)
}
