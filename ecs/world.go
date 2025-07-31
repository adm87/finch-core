package ecs

import (
	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/time"
	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	entities map[EntityID]*Entity

	updateSystems []UpdateSystem
	renderSystems []RenderSystem

	systemsByType map[SystemType]System
}

func NewWorld() *World {
	return &World{
		entities:      make(map[EntityID]*Entity),
		updateSystems: []UpdateSystem{},
		renderSystems: []RenderSystem{},
		systemsByType: make(map[SystemType]System),
	}
}

func (w *World) RegisterSystems(systems ...System) *World {
	for _, system := range systems {
		if system == nil {
			panic(errors.NewNilError("cannot register nil system"))
		}
		if _, exists := w.systemsByType[system.Type()]; exists {
			panic(errors.NewDuplicateError("system already registered"))
		}

		w.systemsByType[system.Type()] = system

		switch sys := system.(type) {
		case UpdateSystem:
			w.updateSystems = append(w.updateSystems, sys)
		case RenderSystem:
			w.renderSystems = append(w.renderSystems, sys)
		default:
			panic(errors.NewInvalidArgumentError("system does not implement UpdateSystem or RenderSystem"))
		}
	}
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
		entities := w.getEntitiesForSystem(system)
		if err := system.Update(entities, t); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) Render(buffer *ebiten.Image, view ebiten.GeoM) error {
	for _, system := range w.renderSystems {
		entities := w.getEntitiesForSystem(system)
		if err := system.Render(entities, buffer, view); err != nil {
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
	w.updateSystems = []UpdateSystem{}
	w.renderSystems = []RenderSystem{}
	w.systemsByType = make(map[SystemType]System)
}
