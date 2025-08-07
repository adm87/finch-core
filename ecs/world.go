package ecs

import (
	"fmt"
	"slices"

	"github.com/adm87/finch-core/errors"
	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	entities           map[EntityID]*Entity
	systemsByType      map[SystemType]System
	fixedUpdateSystems []OrderedSystem[FixedUpdateSystem]
	updateSystems      []OrderedSystem[UpdateSystem]
	renderSystems      []OrderedSystem[RenderSystem]
	view               ebiten.GeoM
}

func NewWorld() *World {
	return &World{
		entities:           make(map[EntityID]*Entity),
		systemsByType:      make(map[SystemType]System),
		fixedUpdateSystems: make([]OrderedSystem[FixedUpdateSystem], 0),
		updateSystems:      make([]OrderedSystem[UpdateSystem], 0),
		renderSystems:      make([]OrderedSystem[RenderSystem], 0),
		view:               ebiten.GeoM{},
	}
}

// RegisterSystems registers one or more systems in the world.
//
// Returns an error if any system is nil, has a nil type, or if a system with the same type is already registered.
//
// Systems are stored by type (UpdateSystem or RenderSystem) and executed in order of their priority (lower values run first).
// If a system implements both UpdateSystem and RenderSystem, the same priority is used for both.
func (w *World) RegisterSystems(systems map[System]int) (*World, error) {
	for system, prio := range systems {
		if system == nil {
			return nil, errors.NewNilError("system cannot be nil")
		}
		if system.Type().IsNil() {
			return nil, errors.NewNilError("system type cannot be nil")
		}
		if _, exists := w.systemsByType[system.Type()]; exists {
			return nil, errors.NewDuplicateError("system already registered: " + fmt.Sprintf("%v", system.Type()))
		}

		fixUpdateSystem, isFixed := system.(FixedUpdateSystem)
		updateSystem, isUpdate := system.(UpdateSystem)
		renderSystem, isRender := system.(RenderSystem)

		if !isFixed && !isUpdate && !isRender {
			return nil, errors.NewInvalidArgumentError("system must implement at least one of FixedUpdateSystem, UpdateSystem or RenderSystem interfaces")
		}

		if isFixed {
			w.fixedUpdateSystems = append(w.fixedUpdateSystems, OrderedSystem[FixedUpdateSystem]{Sys: fixUpdateSystem, Order: prio})
			slices.SortFunc(w.fixedUpdateSystems, func(a, b OrderedSystem[FixedUpdateSystem]) int {
				return a.Order - b.Order
			})
		}
		if isUpdate {
			w.updateSystems = append(w.updateSystems, OrderedSystem[UpdateSystem]{Sys: updateSystem, Order: prio})
			slices.SortFunc(w.updateSystems, func(a, b OrderedSystem[UpdateSystem]) int {
				return a.Order - b.Order
			})
		}
		if isRender {
			w.renderSystems = append(w.renderSystems, OrderedSystem[RenderSystem]{Sys: renderSystem, Order: prio})
			slices.SortFunc(w.renderSystems, func(a, b OrderedSystem[RenderSystem]) int {
				return a.Order - b.Order
			})
		}

		w.systemsByType[system.Type()] = system
	}
	return w, nil
}

func (w *World) AddEntities(entities ...*Entity) (*World, error) {
	for _, entity := range entities {
		if entity == nil {
			return nil, errors.NewNilError("cannot add nil entity to world")
		}
		if _, exists := w.entities[entity.ID()]; exists {
			return nil, errors.NewDuplicateError("entity already exists")
		}
		w.entities[entity.ID()] = entity
	}
	return w, nil
}

func (w *World) RemoveEntities(ids ...EntityID) (*World, error) {
	for _, id := range ids {
		if id.IsNil() {
			return nil, errors.NewNilError("cannot remove entity with nil ID")
		}
		if _, exists := w.entities[id]; !exists {
			continue
		}
		delete(w.entities, id)
	}
	return w, nil
}

func (w *World) GetEntity(id EntityID) (*Entity, bool, error) {
	if id.IsNil() {
		return nil, false, errors.NewNilError("cannot get entity with nil ID")
	}
	entity, exists := w.entities[id]
	return entity, exists, nil
}

func (w *World) GetSystem(systemType SystemType) (System, bool, error) {
	if systemType.IsNil() {
		return nil, false, errors.NewNilError("cannot get system with nil type")
	}
	system, exists := w.systemsByType[systemType]
	return system, exists, nil
}

func (w *World) FixedUpdate() error {
	for _, system := range w.fixedUpdateSystems {
		entities, err := w.internal_filter_entities(system.Sys.Filter())
		if err != nil {
			return err
		}
		if err := system.Sys.FixedUpdate(entities); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) Update() error {
	for _, system := range w.updateSystems {
		entities, err := w.internal_filter_entities(system.Sys.Filter())
		if err != nil {
			return err
		}
		if err := system.Sys.Update(entities); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) Render(screen *ebiten.Image) error {
	for _, system := range w.renderSystems {
		entities, err := w.internal_filter_entities(system.Sys.Filter())
		if err != nil {
			return err
		}
		if err := system.Sys.Render(entities, screen, w.view); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) Clear() {
	w.entities = make(map[EntityID]*Entity)
	w.updateSystems = []OrderedSystem[UpdateSystem]{}
	w.renderSystems = []OrderedSystem[RenderSystem]{}
	w.systemsByType = make(map[SystemType]System)
}

func (w *World) View() ebiten.GeoM {
	return w.view
}

func (w *World) SetView(view ebiten.GeoM) {
	w.view = view
}

func (w *World) internal_filter_entities(filter []ComponentType) ([]*Entity, error) {
	entities := make([]*Entity, 0, len(w.entities))

	for _, entity := range w.entities {
		hasComponent, err := entity.HasComponents(filter...)
		if err != nil {
			return nil, err
		}
		if !hasComponent {
			continue
		}
		entities = append(entities, entity)
	}

	return entities, nil
}
