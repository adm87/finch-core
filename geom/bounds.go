package geom

// Bounded represents any object that has rectangular bounds.
type Bounded interface {
	comparable

	Bounds() Rect64
}
