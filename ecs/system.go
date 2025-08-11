package ecs

import (
	"slices"
	"strconv"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/hash"
	"github.com/hajimehoshi/ebiten/v2"
)

// SystemType is a unique identifier for a system type.
type SystemType hash.Hash

func (t SystemType) IsNil() bool {
	return t == 0
}

func (t SystemType) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func NewSystemType[T System]() SystemType {
	return SystemType(hash.GetHashFromType[T]())
}

// System is an interface that represents a system in the ECS framework.
type System interface {
	// Filter returns the component types that this system operates on.
	Filter() []ComponentType

	Type() SystemType
}

// EarlyUpdateSystem is an interface for systems that need to perform early updates on entities.
//
// Early updates uses a variable delta time and is called before FixedUpdate and LateUpdate.
type EarlyUpdateSystem interface {
	System

	EarlyUpdate(entities hash.HashSet[Entity], deltaSeconds float64) error
}

// FixedUpdateSystem is an interface for systems that need to perform fixed updates on entities. Recommended for physics and similar systems.
//
// Fixed updates uses a fixed delta time and is called after EarlyUpdate and before LateUpdate.
//
// Since fixed updates are called on fixed intervales, it's possible that FixedUpdate is called zero or more times per frame.
type FixedUpdateSystem interface {
	System

	FixedUpdate(entities hash.HashSet[Entity], fixedDeltaSeconds float64) error
}

// LateUpdateSystem is an interface for systems that need to perform late updates on entities.
//
// Late updates uses a variable delta time and is called after FixedUpdate.
type LateUpdateSystem interface {
	System

	LateUpdate(entities hash.HashSet[Entity], deltaSeconds float64) error
}

// RenderSystem is an interface for systems that need to render entities.
type RenderSystem interface {
	System

	Render(entities hash.HashSet[Entity], buffer *ebiten.Image, view ebiten.GeoM) error
}

// OrderedSystem is used for wrapping an instance of a system with an execution priority.
//
// Lower priorities are executed first.
type OrderedSystem[T System] struct {
	Sys   T
	Order int
}

var (
	systemsByType      = make(map[SystemType]System)
	earlyUpdateSystems = make([]OrderedSystem[EarlyUpdateSystem], 0)
	fixedUpdateSystems = make([]OrderedSystem[FixedUpdateSystem], 0)
	lateUpdateSystems  = make([]OrderedSystem[LateUpdateSystem], 0)
	renderSystems      = make([]OrderedSystem[RenderSystem], 0)
)

var (
	ErrNilSystem                 = errors.NewNilError("system cannot be nil")
	ErrSystemAlreadyRegistered   = errors.NewDuplicateError("system already registered")
	ErrInvalidSystemRegistration = errors.NewInvalidArgumentError("system must implement at least one update or render interface")
)

// RegisterSystems registers one or more systems with the ECS framework.
// Provided systems must be unique and cannot share a SystemType.
//
// Systems must be concrete implementations of the System interfaces. They can implement one or more of the following interfaces:
//   - EarlyUpdateSystem
//   - FixedUpdateSystem
//   - LateUpdateSystem
//   - RenderSystem
func RegisterSystems(systems map[System]int) error {
	for sys, order := range systems {
		if sys == nil {
			return ErrNilSystem
		}
		if _, exists := systemsByType[sys.Type()]; exists {
			return ErrSystemAlreadyRegistered
		}

		earlySys, isEarly := sys.(EarlyUpdateSystem)
		fixedSys, isFixed := sys.(FixedUpdateSystem)
		lateSys, isLate := sys.(LateUpdateSystem)
		renderSys, isRender := sys.(RenderSystem)

		if !isEarly && !isFixed && !isLate && !isRender {
			return ErrInvalidSystemRegistration
		}

		if isEarly {
			earlyUpdateSystems = append(earlyUpdateSystems, OrderedSystem[EarlyUpdateSystem]{Sys: earlySys, Order: order})
		}
		if isFixed {
			fixedUpdateSystems = append(fixedUpdateSystems, OrderedSystem[FixedUpdateSystem]{Sys: fixedSys, Order: order})
		}
		if isLate {
			lateUpdateSystems = append(lateUpdateSystems, OrderedSystem[LateUpdateSystem]{Sys: lateSys, Order: order})
		}
		if isRender {
			renderSystems = append(renderSystems, OrderedSystem[RenderSystem]{Sys: renderSys, Order: order})
		}

		systemsByType[sys.Type()] = sys
	}
	slices.SortFunc(earlyUpdateSystems, func(a, b OrderedSystem[EarlyUpdateSystem]) int {
		return a.Order - b.Order
	})
	slices.SortFunc(fixedUpdateSystems, func(a, b OrderedSystem[FixedUpdateSystem]) int {
		return a.Order - b.Order
	})
	slices.SortFunc(lateUpdateSystems, func(a, b OrderedSystem[LateUpdateSystem]) int {
		return a.Order - b.Order
	})
	slices.SortFunc(renderSystems, func(a, b OrderedSystem[RenderSystem]) int {
		return a.Order - b.Order
	})
	return nil
}

// ProcessUpdateSystems processes all update systems in the ECS framework.
func ProcessUpdateSystems(deltaSeconds, fixedDeltaSeconds float64, frameCount int) error {
	for _, sys := range earlyUpdateSystems {
		entities := FilterEntitiesByComponents(sys.Sys.Filter()...)
		if err := sys.Sys.EarlyUpdate(entities, deltaSeconds); err != nil {
			return err
		}
	}

	for frameCount > 0 {
		for _, sys := range fixedUpdateSystems {
			entities := FilterEntitiesByComponents(sys.Sys.Filter()...)
			if err := sys.Sys.FixedUpdate(entities, fixedDeltaSeconds); err != nil {
				return err
			}
		}
		frameCount--
	}

	for _, sys := range lateUpdateSystems {
		entities := FilterEntitiesByComponents(sys.Sys.Filter()...)
		if err := sys.Sys.LateUpdate(entities, deltaSeconds); err != nil {
			return err
		}
	}

	return nil
}

// ProcessRenderSystems processes all render systems in the ECS framework.
func ProcessRenderSystems(buffer *ebiten.Image, view ebiten.GeoM) error {
	for _, sys := range renderSystems {
		entities := FilterEntitiesByComponents(sys.Sys.Filter()...)
		if err := sys.Sys.Render(entities, buffer, view); err != nil {
			return err
		}
	}
	return nil
}
