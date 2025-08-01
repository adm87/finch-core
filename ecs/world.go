package ecs

import (
	"fmt"
	"slices"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/time"
	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	entities      map[EntityID]*Entity
	systemsByType map[SystemType]System
	updateSystems []PrioritizedSystem[UpdateSystem]
	renderSystems []PrioritizedSystem[RenderSystem]
}

func NewWorld() *World {
	return &World{
		entities:      make(map[EntityID]*Entity),
		systemsByType: make(map[SystemType]System),
		updateSystems: []PrioritizedSystem[UpdateSystem]{},
		renderSystems: []PrioritizedSystem[RenderSystem]{},
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
		if sys, ok := system.(UpdateSystem); ok {
			w.updateSystems = append(w.updateSystems, PrioritizedSystem[UpdateSystem]{Sys: sys, Prio: prio})
			slices.SortFunc(w.updateSystems, func(a, b PrioritizedSystem[UpdateSystem]) int {
				return a.Prio - b.Prio
			})
		}
		if sys, ok := system.(RenderSystem); ok {
			w.renderSystems = append(w.renderSystems, PrioritizedSystem[RenderSystem]{Sys: sys, Prio: prio})
			slices.SortFunc(w.renderSystems, func(a, b PrioritizedSystem[RenderSystem]) int {
				return a.Prio - b.Prio
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

func (w *World) Update(t time.Time) error {
	for _, system := range w.updateSystems {
		entities, err := w.getEntitiesForSystem(system.Sys)
		if err != nil {
			return err
		}
		if err := system.Sys.Update(entities, t); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) Render(buffer *ebiten.Image, view ebiten.GeoM) error {
	for _, system := range w.renderSystems {
		entities, err := w.getEntitiesForSystem(system.Sys)
		if err != nil {
			return err
		}
		if err := system.Sys.Render(entities, buffer, view); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) getEntitiesForSystem(system System) ([]Entity, error) {
	filter := system.Filter()
	entities := make([]Entity, 0, len(w.entities))

	for _, entity := range w.entities {
		hasComponent, err := entity.HasComponents(filter...)
		if err != nil {
			return nil, err
		}
		if !hasComponent {
			continue
		}
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (w *World) Clear() {
	w.entities = make(map[EntityID]*Entity)
	w.updateSystems = []PrioritizedSystem[UpdateSystem]{}
	w.renderSystems = []PrioritizedSystem[RenderSystem]{}
	w.systemsByType = make(map[SystemType]System)
}
