package camera

import (
	"github.com/adm87/finch-core/components/transform"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/geometry"
)

var (
	ErrAmbiguousCameras     = errors.NewAmbiguousError("multiple cameras found")
	ErrCameraEntityNotFound = errors.NewNotFoundError("camera entity not found")
)

var CameraComponentType = ecs.NewComponentType[*CameraComponent]()

type CameraComponent struct {
	*transform.TransformComponent

	Zoom       float64
	ZoomFactor float64

	ViewWidth  float64
	ViewHeight float64
}

func (c *CameraComponent) Type() ecs.ComponentType {
	return CameraComponentType
}

func (c *CameraComponent) Viewport() geometry.Rectangle64 {
	position := c.Position()

	width := c.ViewWidth * c.Zoom
	height := c.ViewHeight * c.Zoom

	left := position.X - width/2
	top := position.Y - height/2

	return geometry.Rectangle64{
		X:      left,
		Y:      top,
		Width:  width,
		Height: height,
	}
}

func NewCameraComponent() *CameraComponent {
	return &CameraComponent{
		TransformComponent: transform.NewTransformComponent(),
		Zoom:               1.0,
		ZoomFactor:         0.1,
	}
}

func FindCameraComponent(world *ecs.World) (*CameraComponent, error) {
	entities := world.FilterEntitiesByComponents(CameraComponentType)

	if len(entities) == 0 {
		return nil, ErrCameraEntityNotFound
	}

	if len(entities) > 1 {
		return nil, ErrAmbiguousCameras
	}

	entity, ok := entities.First()
	if !ok {
		return nil, ErrCameraEntityNotFound
	}

	cameraComponent, _, err := ecs.GetComponent[*CameraComponent](world, entity, CameraComponentType)
	return cameraComponent, err
}
