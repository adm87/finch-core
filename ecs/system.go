package ecs

import (
	"strconv"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrSystemTypeMismatch = errors.NewConflictError("system type mismatch")
)

// SystemType is a unique identifier for a system type.
type SystemType uint64

func (t SystemType) IsNil() bool {
	return t == 0
}

func (t SystemType) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func NewSystemType[T System]() SystemType {
	return SystemType(utils.GetHashFromType[T]())
}

func GetSystem[T System](world *World, st SystemType) (T, bool) {
	var zero T

	sys, exists := world.GetSystem(st)

	if !exists {
		return zero, false
	}

	typedSys, ok := sys.(T)

	if !ok {
		return zero, false
	}

	return typedSys, true
}

// System is an interface that represents a system in the ECS framework.
type System interface {
	IsEnabled() bool

	Enable()
	Disable()

	Type() SystemType
}

// EarlyUpdateSystem is an interface for systems that need to perform early updates on entities.
//
// Early updates uses a variable delta time and is called before FixedUpdate and LateUpdate.
type EarlyUpdateSystem interface {
	System

	EarlyUpdate(world *World, deltaSeconds float64) error
}

// FixedUpdateSystem is an interface for systems that need to perform fixed updates on entities. Recommended for physics and similar systems.
//
// Fixed updates uses a fixed delta time and is called after EarlyUpdate and before LateUpdate.
//
// Since fixed updates are called on fixed intervales, it's possible that FixedUpdate is called zero or more times per frame.
type FixedUpdateSystem interface {
	System

	FixedUpdate(world *World, fixedDeltaSeconds float64) error
}

// LateUpdateSystem is an interface for systems that need to perform late updates on entities.
//
// Late updates uses a variable delta time and is called after FixedUpdate.
type LateUpdateSystem interface {
	System

	LateUpdate(world *World, deltaSeconds float64) error
}

// RenderSystem is an interface for systems that need to render entities.
type RenderSystem interface {
	System

	Render(world *World, buffer *ebiten.Image) error
}

// OrderedSystem is used for wrapping an instance of a system with an execution priority.
//
// Lower priorities are executed first.
type OrderedSystem[T System] struct {
	Sys   T
	Order int
}
