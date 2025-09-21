package geom

type Point64 struct {
	X float64
	Y float64
}

func NewPoint64(x, y float64) Point64 {
	return Point64{X: x, Y: y}
}

func (p Point64) Add(o Point64) Point64 {
	return Point64{X: p.X + o.X, Y: p.Y + o.Y}
}

func (p Point64) Sub(o Point64) Point64 {
	return Point64{X: p.X - o.X, Y: p.Y - o.Y}
}

func (p Point64) Mul(scalar float64) Point64 {
	return Point64{X: p.X * scalar, Y: p.Y * scalar}
}

func (p Point64) Div(scalar float64) Point64 {
	return Point64{X: p.X / scalar, Y: p.Y / scalar}
}

func (p Point64) Equal(o Point64) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p Point64) Length() float64 {
	return (p.X*p.X + p.Y*p.Y)
}

func (p Point64) Distance(o Point64) float64 {
	dx := p.X - o.X
	dy := p.Y - o.Y
	return (dx*dx + dy*dy)
}

func (p Point64) Normalized() Point64 {
	len := p.Length()
	if len == 0 {
		return Point64{X: 0, Y: 0}
	}
	return Point64{X: p.X / len, Y: p.Y / len}
}
