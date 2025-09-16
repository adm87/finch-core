package geom

type Rect64 struct {
	minx, miny, maxx, maxy float64
}

func NewRect64(x, y, w, h float64) Rect64 {
	return Rect64{
		minx: x,
		miny: y,
		maxx: x + w,
		maxy: y + h,
	}
}

func (r Rect64) Min() (x, y float64) {
	return r.minx, r.miny
}

func (r Rect64) Max() (x, y float64) {
	return r.maxx, r.maxy
}

func (r Rect64) Center() (x, y float64) {
	return (r.minx + r.maxx) / 2, (r.miny + r.maxy) / 2
}

func (r Rect64) Width() float64 {
	return r.maxx - r.minx
}

func (r Rect64) Height() float64 {
	return r.maxy - r.miny
}

func (r Rect64) IsEmpty() bool {
	return r.minx >= r.maxx || r.miny >= r.maxy
}

func (r Rect64) Contains(x, y float64) bool {
	return x >= r.minx && x < r.maxx && y >= r.miny && y < r.maxy
}

func (r Rect64) Intersect(o Rect64) (ir Rect64, ok bool) {
	ir.minx = max(r.minx, o.minx)
	ir.miny = max(r.miny, o.miny)
	ir.maxx = min(r.maxx, o.maxx)
	ir.maxy = min(r.maxy, o.maxy)
	if ir.IsEmpty() {
		return ir, false
	}
	return ir, true
}

func (r Rect64) Union(o Rect64) Rect64 {
	if r.IsEmpty() {
		return o
	}
	if o.IsEmpty() {
		return r
	}
	return Rect64{
		minx: min(r.minx, o.minx),
		miny: min(r.miny, o.miny),
		maxx: max(r.maxx, o.maxx),
		maxy: max(r.maxy, o.maxy),
	}
}

func (r Rect64) Inset(dx, dy float64) Rect64 {
	return Rect64{
		minx: r.minx + dx,
		miny: r.miny + dy,
		maxx: r.maxx - dx,
		maxy: r.maxy - dy,
	}
}

func (r Rect64) Offset(dx, dy float64) Rect64 {
	return Rect64{
		minx: r.minx + dx,
		miny: r.miny + dy,
		maxx: r.maxx + dx,
		maxy: r.maxy + dy,
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
