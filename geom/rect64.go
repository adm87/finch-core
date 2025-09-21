package geom

type Rect64 struct {
	x, y, w, h float64
}

func NewRect64(x, y, w, h float64) Rect64 {
	return Rect64{x, y, w, h}
}

func (r *Rect64) SetXY(x, y float64) {
	r.x = x
	r.y = y
}

func (r *Rect64) SetSize(w, h float64) {
	r.w = w
	r.h = h
}

func (r Rect64) Min() (float64, float64) {
	return r.x, r.y
}

func (r Rect64) Max() (float64, float64) {
	return r.x + r.w, r.y + r.h
}

func (r Rect64) Width() float64 {
	return r.w
}

func (r Rect64) Height() float64 {
	return r.h
}

func (r Rect64) ContainsXY(x, y float64) bool {
	return x >= r.x && x < r.x+r.w && y >= r.y && y < r.y+r.h
}
