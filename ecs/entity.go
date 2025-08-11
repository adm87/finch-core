package ecs

import (
	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/hash"
	"github.com/google/uuid"
)

// Entity represents a unique entity in the ECS framework.
type Entity uuid.UUID

func (id Entity) IsNil() bool {
	return id == Entity(uuid.Nil)
}

func (id Entity) String() string {
	return uuid.UUID(id).String()
}

var (
	knownEntities            = make(hash.HashSet[Entity])
	entitiesByComponentTypes = make(map[ComponentType]hash.HashSet[Entity])
	componentsByEntity       = make(map[Entity]map[ComponentType]Component)
)

var (
	ErrNilEntity              = errors.NewNilError("entity cannot be nil")
	ErrNilComponent           = errors.NewNilError("component cannot be nil")
	ErrUnknownEntity          = errors.NewNotFoundError("unknown entity")
	ErrComponentNotFound      = errors.NewNotFoundError("component not found for entity")
	ErrEntityAlreadyExists    = errors.NewDuplicateError("entity already exists")
	ErrComponentAlreadyExists = errors.NewDuplicateError("component already exists for entity")
	ErrComponentTypeMismatch  = errors.NewConflictError("component type mismatch")
)

var NilEntity = Entity(uuid.Nil)

// NewEntity creates a new entity and adds it to the known entities.
func NewEntity() Entity {
	id := Entity(uuid.New())
	knownEntities.Add(id)
	return id
}

// NewEntityWithComponents like NewEntity but assigns the provided components to the new entity.
//
// If adding components fails, the entity will be removed from the known entities.
func NewEntityWithComponents(components ...Component) (Entity, error) {
	entity := NewEntity()

	if len(components) == 0 {
		return entity, nil
	}

	if err := AddComponents(entity, components...); err != nil {
		knownEntities.Remove(entity)
		return NilEntity, err
	}

	return entity, nil
}

// AddEntity adds an existing Entity to the known entities
//
// Use this when deserializing entities.
func AddEntity(entity Entity) error {
	if entity.IsNil() {
		return ErrNilEntity
	}

	if knownEntities.Contains(entity) {
		return ErrEntityAlreadyExists
	}

	knownEntities.Add(entity)
	return nil
}

// AddEntityWithComponents like AddEntity but also assigns the provided components to the entity.
// If adding components fails, the entity will be removed from the known entities.
//
// Use this when deserializing entities with components.
func AddEntityWithComponents(entity Entity, components ...Component) error {
	if err := AddEntity(entity); err != nil {
		return err
	}

	if err := AddComponents(entity, components...); err != nil {
		knownEntities.Remove(entity)
		return err
	}

	return nil
}

// AddComponents assigns one or more components to an entity.
// Components must be unique concrete types implementing the ecs.Component interface.
//
// Example:
//
//	AddComponents(entity, &TransformComponent{}, &RenderComponent{})
func AddComponents(entity Entity, components ...Component) error {
	if entity.IsNil() {
		return ErrNilEntity
	}

	if !knownEntities.Contains(entity) {
		return ErrUnknownEntity
	}

	for _, component := range components {
		if component == nil {
			return ErrNilComponent
		}

		ct := component.Type()

		if ct.IsNil() {
			return ErrNilComponent
		}

		if _, ok := entitiesByComponentTypes[ct]; !ok {
			entitiesByComponentTypes[ct] = make(hash.HashSet[Entity])
		}

		if _, exists := entitiesByComponentTypes[ct][entity]; exists {
			return ErrComponentAlreadyExists
		}

		entitiesByComponentTypes[ct].Add(entity)

		if _, exists := componentsByEntity[entity]; !exists {
			componentsByEntity[entity] = make(map[ComponentType]Component)
		}

		if _, exists := componentsByEntity[entity][ct]; exists {
			return ErrComponentAlreadyExists // We should never get here
		}

		componentsByEntity[entity][ct] = component
	}

	return nil
}

// RemoveComponents unassigns one or more components from an entity.
//
// Removing a component will dispose of it and should be used considered unusable.
func RemoveComponents(entity Entity, componentTypes ...ComponentType) error {
	if entity.IsNil() {
		return ErrNilEntity
	}

	if !knownEntities.Contains(entity) {
		return ErrUnknownEntity
	}

	for _, ct := range componentTypes {
		if ct.IsNil() {
			return ErrNilComponent
		}

		if _, exists := entitiesByComponentTypes[ct]; !exists {
			continue
		}

		delete(entitiesByComponentTypes[ct], entity)
		if len(entitiesByComponentTypes[ct]) == 0 {
			delete(entitiesByComponentTypes, ct)
		}

		if _, exists := componentsByEntity[entity]; !exists {
			continue
		}

		component, exists := componentsByEntity[entity][ct]

		if !exists {
			continue
		}

		component.Dispose()

		delete(componentsByEntity[entity], ct)
		if len(componentsByEntity[entity]) == 0 {
			delete(componentsByEntity, entity)
		}
	}

	return nil
}

// GetComponent returns an entity's assigned component of a specific type.
// T must be a concrete ecs.Component type pointer, and must match the provided componentType.
//
// Example:
//
//	GetComponent[*TransformComponent](entity, TransformComponentType)
func GetComponent[T Component](entity Entity, componentType ComponentType) (T, bool, error) {
	var component T

	if entity.IsNil() {
		return component, false, ErrNilEntity
	}

	if !knownEntities.Contains(entity) {
		return component, false, ErrUnknownEntity
	}

	if _, exists := entitiesByComponentTypes[componentType]; !exists || !entitiesByComponentTypes[componentType].Contains(entity) {
		return component, false, nil
	}

	if _, exists := componentsByEntity[entity]; !exists || componentsByEntity[entity][componentType] == nil {
		return component, false, nil
	}

	component, ok := componentsByEntity[entity][componentType].(T)
	if !ok {
		return component, false, ErrComponentTypeMismatch
	}

	return component, true, nil
}

// FilterEntitiesByComponents returns a set of entities that have all of the specified component types.
func FilterEntitiesByComponents(componentTypes ...ComponentType) hash.HashSet[Entity] {
	if len(componentTypes) == 0 {
		return knownEntities
	}

	entities := make([]hash.HashSet[Entity], 0)

	for _, ct := range componentTypes {
		if ct.IsNil() {
			return hash.HashSet[Entity]{}
		}

		set, ok := entitiesByComponentTypes[ct]

		if !ok {
			return hash.HashSet[Entity]{}
		}

		entities = append(entities, set)
	}

	return hash.IntersectHashSets(entities...)
}
