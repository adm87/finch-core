package geometry

import (
	"fmt"
	"math"
)

// =================================================================
// Point
// =================================================================

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

// =================================================================
// Point64
// =================================================================

type Point64 struct {
	X, Y float64
}

// Add adds another Point64 to this Point64.
func (p *Point64) Add(other Point64) {
	p.X += other.X
	p.Y += other.Y
}

// Subtract subtracts another Point64 from this Point64.
func (p *Point64) Subtract(other Point64) {
	p.X -= other.X
	p.Y -= other.Y
}

// Scale scales the Point64 by a scalar value.
func (p *Point64) Scale(scalar float64) {
	p.X *= scalar
	p.Y *= scalar
}

// DistanceTo calculates the distance from this Point64 to another Point64.
func (p Point64) DistanceTo(other Point64) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Length calculates the length of the Point64 vector from the origin (0, 0).
func (p Point64) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Normalize returns a normalized version of the Point64 vector.
func (p Point64) Normalize() Point64 {
	length := p.Length()
	if length == 0 {
		return Point64{0, 0}
	}
	return Point64{p.X / length, p.Y / length}
}

func (p Point64) String() string {
	return fmt.Sprintf("Point64(X: %.2f, Y: %.2f)", p.X, p.Y)
}
