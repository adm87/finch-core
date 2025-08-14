package ecs

import (
	"slices"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/hash"
	"github.com/adm87/finch-core/linq"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrNilEntity                 = errors.NewNilError("entity cannot be nil")
	ErrNilComponent              = errors.NewNilError("component cannot be nil")
	ErrNilSystem                 = errors.NewNilError("system cannot be nil")
	ErrUnknownEntity             = errors.NewNotFoundError("unknown entity")
	ErrComponentNotFound         = errors.NewNotFoundError("component not found for entity")
	ErrEntityAlreadyExists       = errors.NewDuplicateError("entity already exists")
	ErrComponentAlreadyExists    = errors.NewDuplicateError("component already exists for entity")
	ErrSystemAlreadyRegistered   = errors.NewDuplicateError("system already registered")
	ErrComponentTypeMismatch     = errors.NewConflictError("component type mismatch")
	ErrInvalidSystemRegistration = errors.NewInvalidArgumentError("system must implement at least one update or render interface")
)

// ECSWorld is set of registered systems, entities, and components
type ECSWorld struct {
	entities                hash.HashSet[Entity]
	entitiesByComponentType map[ComponentType]hash.HashSet[Entity]
	componentsByEntity      map[Entity]map[ComponentType]Component
	systems                 map[SystemType]System
	earlyUpdates            []OrderedSystem[EarlyUpdateSystem]
	fixedUpdates            []OrderedSystem[FixedUpdateSystem]
	lateUpdates             []OrderedSystem[LateUpdateSystem]
	renderSystems           []OrderedSystem[RenderSystem]
}

func NewWorld() *ECSWorld {
	return &ECSWorld{
		entities:                make(hash.HashSet[Entity]),
		entitiesByComponentType: make(map[ComponentType]hash.HashSet[Entity]),
		componentsByEntity:      make(map[Entity]map[ComponentType]Component),
		systems:                 make(map[SystemType]System),
		earlyUpdates:            make([]OrderedSystem[EarlyUpdateSystem], 0),
		fixedUpdates:            make([]OrderedSystem[FixedUpdateSystem], 0),
		lateUpdates:             make([]OrderedSystem[LateUpdateSystem], 0),
		renderSystems:           make([]OrderedSystem[RenderSystem], 0),
	}
}

// =================================================================
// Entity Management
// =================================================================

// NewEntity creates a new entity and adds it to the world.
func (w *ECSWorld) NewEntity() Entity {
	id := Entity(uuid.New())
	w.entities.Add(id)
	return id
}

// NewEntityWithComponents creates a new entity with the given components and adds it to the world.
//
// If any of the components fail to be added, the entity will be removed from the world and a NilEntity is returned along with the error.
func (w *ECSWorld) NewEntityWithComponents(components ...Component) (Entity, error) {
	entity := w.NewEntity()

	if err := w.AddComponents(entity, components...); err != nil {
		w.entities.Remove(entity)
		return NilEntity, err
	}

	return entity, nil
}

// AddEntity adds an existing Entity to the world.
//
// Useful when deserializing entities.
func (w *ECSWorld) AddEntity(entity Entity) error {
	if entity.IsNil() {
		return ErrNilEntity
	}

	if w.entities.Contains(entity) {
		return ErrEntityAlreadyExists
	}

	w.entities.Add(entity)
	return nil
}

// AddEntityWithComponents adds an existing Entity with the given components to the world.
//
// If any of the components fail to be added, the entity will be removed from the world and the error is returned.
func (w *ECSWorld) AddEntityWithComponents(entity Entity, components ...Component) error {
	if err := w.AddEntity(entity); err != nil {
		return err
	}

	if err := w.AddComponents(entity, components...); err != nil {
		w.entities.Remove(entity)
		return err
	}

	return nil
}

// Removes an entity and its components from the world.
//
// Removing an entity will make it and its components unusable.
//
// If a component fails to be removed, the entity will still be removed from the world.
func (w *ECSWorld) RemoveEntity(entity Entity) (bool, error) {
	if entity.IsNil() {
		return false, ErrNilEntity
	}

	if !w.entities.Contains(entity) {
		return false, ErrUnknownEntity
	}

	w.entities.Remove(entity)

	cts := linq.Keys(w.componentsByEntity[entity])
	for _, ct := range cts {
		if err := w.RemoveComponent(entity, ct); err != nil {
			return false, err
		}
	}

	return true, nil
}

// FilterEntitiesByComponents returns a set of entities that have all of the specified component types.
func (w *ECSWorld) FilterEntitiesByComponents(componentTypes ...ComponentType) hash.HashSet[Entity] {
	if len(componentTypes) == 0 {
		return w.entities
	}

	sets := make([]hash.HashSet[Entity], 0, len(componentTypes))
	for _, ct := range componentTypes {
		if ct.IsNil() {
			return hash.HashSet[Entity]{}
		}
		set, ok := w.entitiesByComponentType[ct]
		if !ok {
			return hash.HashSet[Entity]{}
		}
		sets = append(sets, set)
	}

	return hash.IntersectHashSets(sets...)
}

// =================================================================
// Component Management
// =================================================================

// AddComponent adds a component to an entity in the world.
func (w *ECSWorld) AddComponent(entity Entity, component Component) error {
	if entity.IsNil() {
		return ErrNilEntity
	}

	if !w.entities.Contains(entity) {
		return ErrUnknownEntity
	}

	if component == nil {
		return ErrNilComponent
	}

	ct := component.Type()

	if _, exists := w.entitiesByComponentType[ct]; !exists {
		w.entitiesByComponentType[ct] = make(hash.HashSet[Entity])
	}

	if _, exists := w.entitiesByComponentType[ct][entity]; exists {
		return ErrComponentAlreadyExists
	}

	w.entitiesByComponentType[ct].Add(entity)

	if _, exists := w.componentsByEntity[entity]; !exists {
		w.componentsByEntity[entity] = make(map[ComponentType]Component)
	}

	w.componentsByEntity[entity][ct] = component

	return nil
}

// AddComponents adds multiple components to an entity in the world.
func (w *ECSWorld) AddComponents(entity Entity, components ...Component) error {
	for _, c := range components {
		if err := w.AddComponent(entity, c); err != nil {
			return err
		}
	}
	return nil
}

// RemoveComponent removes a component from an entity in the world.
func (w *ECSWorld) RemoveComponent(entity Entity, ct ComponentType) error {
	if entity.IsNil() {
		return ErrNilEntity
	}

	if !w.entities.Contains(entity) {
		return ErrUnknownEntity
	}

	if ct.IsNil() {
		return ErrNilComponent
	}

	if _, exists := w.entitiesByComponentType[ct]; !exists {
		return ErrComponentNotFound
	}

	if _, exists := w.entitiesByComponentType[ct][entity]; !exists {
		return ErrComponentNotFound
	}

	w.entitiesByComponentType[ct].Remove(entity)
	if len(w.entitiesByComponentType[ct]) == 0 {
		delete(w.entitiesByComponentType, ct)
	}

	// This should be safe, we don't map the component in the Add methods unless it's a valid addition.
	if component := w.componentsByEntity[entity][ct]; component != nil {
		component.Dispose()

		delete(w.componentsByEntity[entity], ct)
		if len(w.componentsByEntity[entity]) == 0 {
			delete(w.componentsByEntity, entity)
		}
	}

	return nil
}

// RemoveComponents removes multiple components from an entity in the world.
func (w *ECSWorld) RemoveComponents(entity Entity, cts ...ComponentType) error {
	for _, ct := range cts {
		if err := w.RemoveComponent(entity, ct); err != nil {
			return err
		}
	}
	return nil
}

// GetComponent returns a component from an entity within the world.
func (w *ECSWorld) GetComponent(entity Entity, ct ComponentType) (Component, bool, error) {
	if entity.IsNil() {
		return nil, false, ErrNilEntity
	}

	if !w.entities.Contains(entity) {
		return nil, false, ErrUnknownEntity
	}

	if ct.IsNil() {
		return nil, false, ErrNilComponent
	}

	if _, exists := w.entitiesByComponentType[ct]; !exists {
		return nil, false, nil
	}

	if _, exists := w.entitiesByComponentType[ct][entity]; !exists {
		return nil, false, nil
	}

	component, exists := w.componentsByEntity[entity][ct]
	if !exists {
		return nil, false, nil
	}

	return component, true, nil
}

// GetComponents returns multiple components from an entity within the world.
func (w *ECSWorld) GetComponents(entity Entity, cts ...ComponentType) (map[ComponentType]Component, error) {
	components := make(map[ComponentType]Component)

	for _, ct := range cts {
		if component, found, err := w.GetComponent(entity, ct); err != nil {
			return nil, err
		} else if found {
			components[ct] = component
		}
	}

	return components, nil
}

// =================================================================
// System Management
// =================================================================

func (w *ECSWorld) RegisterSystems(systems map[System]int) error {
	for sys, order := range systems {
		if sys == nil {
			return ErrNilSystem
		}

		st := sys.Type()

		if st.IsNil() {
			return ErrNilSystem
		}

		if _, exists := w.systems[st]; exists {
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
			w.earlyUpdates = append(w.earlyUpdates, OrderedSystem[EarlyUpdateSystem]{Sys: earlySys, Order: order})
		}
		if isFixed {
			w.fixedUpdates = append(w.fixedUpdates, OrderedSystem[FixedUpdateSystem]{Sys: fixedSys, Order: order})
		}
		if isLate {
			w.lateUpdates = append(w.lateUpdates, OrderedSystem[LateUpdateSystem]{Sys: lateSys, Order: order})
		}
		if isRender {
			w.renderSystems = append(w.renderSystems, OrderedSystem[RenderSystem]{Sys: renderSys, Order: order})
		}

		w.systems[sys.Type()] = sys
	}
	slices.SortFunc(w.earlyUpdates, func(a, b OrderedSystem[EarlyUpdateSystem]) int {
		return a.Order - b.Order
	})
	slices.SortFunc(w.fixedUpdates, func(a, b OrderedSystem[FixedUpdateSystem]) int {
		return a.Order - b.Order
	})
	slices.SortFunc(w.lateUpdates, func(a, b OrderedSystem[LateUpdateSystem]) int {
		return a.Order - b.Order
	})
	slices.SortFunc(w.renderSystems, func(a, b OrderedSystem[RenderSystem]) int {
		return a.Order - b.Order
	})
	return nil
}

func (w *ECSWorld) ProcessUpdateSystems(deltaSeconds, fixedDeltaSeconds float64, frameCount int) error {
	for _, sys := range w.earlyUpdates {
		if err := sys.Sys.EarlyUpdate(w, deltaSeconds); err != nil {
			return err
		}
	}

	for frameCount > 0 {
		for _, sys := range w.fixedUpdates {
			if err := sys.Sys.FixedUpdate(w, fixedDeltaSeconds); err != nil {
				return err
			}
		}
		frameCount--
	}

	for _, sys := range w.lateUpdates {
		if err := sys.Sys.LateUpdate(w, deltaSeconds); err != nil {
			return err
		}
	}

	return nil
}

func (w *ECSWorld) ProcessRenderSystems(screen *ebiten.Image) error {
	for _, sys := range w.renderSystems {
		if err := sys.Sys.Render(w, screen); err != nil {
			return err
		}
	}
	return nil
}
