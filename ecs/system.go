package ecs

import (
	"hash"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// SystemType is a unique identifier for a system type.
type SystemType hash.Hash

// System is an interface that represents a system in the ECS framework.
type System interface {
	// Filter returns the component types that this system operates on.
	Filter() []ComponentType

	Type() SystemType
}

// UpdateSystem is an interface for systems that need to update entities.
type UpdateSystem interface {
	System

	Update(entities []Entity, t time.Time) error
}

// RenderSystem is an interface for systems that need to render entities.
type RenderSystem interface {
	System

	Render(entities []Entity, buffer *ebiten.Image, view ebiten.GeoM) error
}
