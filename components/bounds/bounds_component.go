package bounds

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
)

var BoundsComponentType = ecs.NewComponentType[*BoundsComponent]()

type BoundsComponent struct {
	Size   geometry.Point64
	Anchor geometry.Point64
}

func NewBoundsComponent(size, anchor geometry.Point64) *BoundsComponent {
	return &BoundsComponent{
		Size:   size,
		Anchor: anchor,
	}
}

func (b *BoundsComponent) Type() ecs.ComponentType {
	return BoundsComponentType
}

func (b *BoundsComponent) Left() float64 {
	return -b.Size.X * b.Anchor.X
}

func (b *BoundsComponent) Right() float64 {
	return b.Size.X * (1 - b.Anchor.X)
}

func (b *BoundsComponent) Top() float64 {
	return -b.Size.Y * b.Anchor.Y
}

func (b *BoundsComponent) Bottom() float64 {
	return b.Size.Y * (1 - b.Anchor.Y)
}

func (b *BoundsComponent) Center() (x, y float64) {
	x = b.Left() + b.Size.X/2
	y = b.Top() + b.Size.Y/2
	return
}

func (b *BoundsComponent) Min() (x, y float64) {
	x = b.Left()
	y = b.Top()
	return
}

func (b *BoundsComponent) Max() (x, y float64) {
	x = b.Right()
	y = b.Bottom()
	return
}

func (b *BoundsComponent) AABB(worldPosition geometry.Point64) geometry.Rectangle64 {
	return geometry.Rectangle64{
		X:      b.Left() + worldPosition.X,
		Y:      b.Top() + worldPosition.Y,
		Width:  b.Size.X,
		Height: b.Size.Y,
	}
}
