package transform

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-core/hash"
	"github.com/adm87/finch-core/math"
	"github.com/hajimehoshi/ebiten/v2"
)

var TransformComponentType = ecs.ComponentType(hash.GetHashFromType[TransformComponent]())

type TransformComponent struct {
	position geometry.Point64
	scale    geometry.Point64
	origin   geometry.Point64

	deg, rad float64

	localMatrix ebiten.GeoM
	worldMatrix ebiten.GeoM

	localDirty bool
	worldDirty bool
}

func NewTransformComponent() *TransformComponent {
	return NewTransformComponentWith(
		geometry.Point64{X: 0, Y: 0},
		geometry.Point64{X: 1, Y: 1},
		geometry.Point64{X: 0, Y: 0},
		0,
	)
}

func NewTransformComponentWith(position geometry.Point64, scale geometry.Point64, origin geometry.Point64, deg float64) *TransformComponent {
	return &TransformComponent{
		position:    position,
		scale:       scale,
		origin:      origin,
		deg:         deg,
		rad:         deg * math.DegToRad,
		localMatrix: ebiten.GeoM{},
		worldMatrix: ebiten.GeoM{},
		localDirty:  true,
		worldDirty:  true,
	}
}

func (t *TransformComponent) Type() ecs.ComponentType {
	return TransformComponentType
}

func (t *TransformComponent) Position() geometry.Point64 {
	return t.position
}

func (t *TransformComponent) SetPosition(position geometry.Point64) {
	if t.position.X == position.X && t.position.Y == position.Y {
		return
	}
	t.position = position
	t.localDirty = true
	t.worldDirty = true
}

func (t *TransformComponent) Scale() geometry.Point64 {
	return t.scale
}

func (t *TransformComponent) SetScale(scale geometry.Point64) {
	if t.scale.X == scale.X && t.scale.Y == scale.Y {
		return
	}
	t.scale = scale
	t.localDirty = true
	t.worldDirty = true
}

func (t *TransformComponent) Origin() geometry.Point64 {
	return t.origin
}

func (t *TransformComponent) SetOrigin(origin geometry.Point64) {
	if t.origin.X == origin.X && t.origin.Y == origin.Y {
		return
	}
	t.origin = origin
	t.localDirty = true
	t.worldDirty = true
}

func (t *TransformComponent) Rotation() (float64, float64) {
	return t.deg, t.rad
}

func (t *TransformComponent) SetRotation(deg float64) {
	if deg < 0 {
		deg += 360
	}
	if deg >= 360 {
		deg -= 360
	}
	if deg == t.deg {
		return
	}
	t.deg = deg
	t.rad = deg * math.DegToRad
	t.localDirty = true
	t.worldDirty = true
}

func (t *TransformComponent) LocalMatrix() ebiten.GeoM {
	if t.localDirty {
		t.localMatrix.Reset()
		t.localMatrix.Translate(-t.origin.X, -t.origin.Y)
		t.localMatrix.Scale(t.scale.X, t.scale.Y)
		t.localMatrix.Rotate(t.rad)
		t.localMatrix.Translate(t.position.X, t.position.Y)
		t.localDirty = false
	}
	return t.localMatrix
}

func (t *TransformComponent) WorldMatrix() ebiten.GeoM {
	if t.worldDirty {
		t.worldMatrix.Reset()
		t.worldMatrix.Concat(t.LocalMatrix())
		t.worldDirty = false
	}
	return t.worldMatrix
}
