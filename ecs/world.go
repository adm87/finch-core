package ecs

import (
	"slices"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/time"
	"github.com/hajimehoshi/ebiten/v2"
)

type WorldSys[T System] struct {
	Sys  T
	Prio int
}

type World struct {
	entities      map[EntityID]*Entity
	systemsByType map[SystemType]System
	updateSystems []WorldSys[UpdateSystem]
	renderSystems []WorldSys[RenderSystem]
}

func NewWorld() *World {
	return &World{
		entities:      make(map[EntityID]*Entity),
		systemsByType: make(map[SystemType]System),
		updateSystems: []WorldSys[UpdateSystem]{},
		renderSystems: []WorldSys[RenderSystem]{},
	}
}

// RegisterSystem registers a system in the world.
//
// Panics if the system is nil, if the system type is already registered, or if the system does not implement UpdateSystem or RenderSystem.
//
// Systems will be stored based on their type (UpdateSystem or RenderSystem) and will be executed in the order of their registered priority.
//
// TODO: Currently if a system is both an UpdateSystem and RenderSystem, it uses the same priority for both categorizations. This might be addressed later.
func (w *World) RegisterSystem(system System, priority int) *World {
	if system == nil {
		panic(errors.NewNilError("cannot register nil system"))
	}
	if _, exists := w.systemsByType[system.Type()]; exists {
		panic(errors.NewDuplicateError("system already registered"))
	}
	if sys := system.(UpdateSystem); sys != nil {
		w.updateSystems = append(w.updateSystems, WorldSys[UpdateSystem]{Sys: sys, Prio: priority})
		slices.SortFunc(w.updateSystems, func(a, b WorldSys[UpdateSystem]) int {
			return a.Prio - b.Prio
		})
	}
	if sys := system.(RenderSystem); sys != nil {
		w.renderSystems = append(w.renderSystems, WorldSys[RenderSystem]{Sys: sys, Prio: priority})
		slices.SortFunc(w.renderSystems, func(a, b WorldSys[RenderSystem]) int {
			return a.Prio - b.Prio
		})
	}
	w.systemsByType[system.Type()] = system
	return w
}

func (w *World) AddEntities(entities ...*Entity) *World {
	for _, entity := range entities {
		if entity == nil {
			panic(errors.NewNilError("cannot add nil entity to world"))
		}
		if _, exists := w.entities[entity.ID()]; exists {
			panic(errors.NewDuplicateError("entity already exists"))
		}
		w.entities[entity.ID()] = entity
	}
	return w
}

func (w *World) RemoveEntities(ids ...EntityID) *World {
	for _, id := range ids {
		if id.IsNil() {
			panic(errors.NewNilError("cannot remove entity with nil ID"))
		}
		if _, exists := w.entities[id]; !exists {
			continue
		}
		delete(w.entities, id)
	}
	return w
}

func (w *World) GetEntity(id EntityID) (*Entity, bool) {
	if id.IsNil() {
		panic(errors.NewNilError("cannot get entity with nil ID"))
	}
	entity, exists := w.entities[id]
	return entity, exists
}

func (w *World) GetSystem(systemType SystemType) (System, bool) {
	if systemType.IsNil() {
		panic(errors.NewNilError("cannot get system with nil type"))
	}
	system, exists := w.systemsByType[systemType]
	return system, exists
}

func (w *World) Update(t time.Time) error {
	for _, system := range w.updateSystems {
		entities := w.getEntitiesForSystem(system.Sys)
		if err := system.Sys.Update(entities, t); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) Render(buffer *ebiten.Image, view ebiten.GeoM) error {
	for _, system := range w.renderSystems {
		entities := w.getEntitiesForSystem(system.Sys)
		if err := system.Sys.Render(entities, buffer, view); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) getEntitiesForSystem(system System) []Entity {
	filter := system.Filter()
	entities := make([]Entity, 0, len(w.entities))

	for _, entity := range w.entities {
		if entity.HasComponents(filter...) {
			entities = append(entities, *entity)
		}
	}

	return entities
}

func (w *World) Clear() {
	w.entities = make(map[EntityID]*Entity)
	w.updateSystems = []WorldSys[UpdateSystem]{}
	w.renderSystems = []WorldSys[RenderSystem]{}
	w.systemsByType = make(map[SystemType]System)
}
