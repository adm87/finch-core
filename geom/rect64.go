package geom

type Rect64 struct {
	X, Y, Width, Height float64
}

func NewRect64(x, y, w, h float64) Rect64 {
	return Rect64{X: x, Y: y, Width: w, Height: h}
}

func (r *Rect64) SetXY(x, y float64) {
	r.X = x
	r.Y = y
}

func (r *Rect64) SetSize(w, h float64) {
	r.Width = w
	r.Height = h
}

func (r Rect64) Min() (float64, float64) {
	return r.X, r.Y
}

func (r Rect64) Max() (float64, float64) {
	return r.X + r.Width, r.Y + r.Height
}

func (r Rect64) ContainsXY(x, y float64) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

func (r Rect64) Intersects(o Rect64) bool {
	return r.X < o.X+o.Width && r.X+r.Width > o.X && r.Y < o.Y+o.Height && r.Y+r.Height > o.Y
}
