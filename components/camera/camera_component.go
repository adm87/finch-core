package camera

import (
	"github.com/adm87/finch-core/components/transform"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/errors"
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
}

func (c *CameraComponent) Type() ecs.ComponentType {
	return CameraComponentType
}

func (c *CameraComponent) Dispose() {
	c.TransformComponent = nil
}

func NewCameraComponent() *CameraComponent {
	return &CameraComponent{
		TransformComponent: transform.NewTransformComponent(),
		Zoom:               1.0,
		ZoomFactor:         0.1,
	}
}

func FindCameraComponent(world *ecs.ECSWorld) (*CameraComponent, error) {
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
