package ecs

import (
	"github.com/adm87/finch-core/hash"
	"github.com/hajimehoshi/ebiten/v2"
)

// PrioritizedSystem is used for wrapping an instance of a system with an execution priority.
//
// Lower priorities are executed first.
type PrioritizedSystem[T System] struct {
	Sys  T
	Prio int
}

// SystemType is a unique identifier for a system type.
type SystemType hash.Hash

func (t SystemType) IsNil() bool {
	return t == 0
}

// System is an interface that represents a system in the ECS framework.
type System interface {
	// Filter returns the component types that this system operates on.
	Filter() []ComponentType

	Type() SystemType
}

// UpdateSystem is an interface for systems that need to update entities.
type UpdateSystem interface {
	System

	FixedUpdate(entities []*Entity) error

	Update(entities []*Entity) error
}

// RenderSystem is an interface for systems that need to render entities.
type RenderSystem interface {
	System

	Render(entities []*Entity, buffer *ebiten.Image, view ebiten.GeoM) error
}
