package ecs

import (
	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/hash"
	"github.com/google/uuid"
)

// EntityID is a unique identifier for an entity in the ECS framework.
type EntityID uuid.UUID

func (e EntityID) IsNil() bool {
	return e == EntityID(uuid.Nil)
}

type Entity struct {
	id         EntityID
	components map[ComponentType]Component
	tags       hash.HashSet[string]
}

func NewEntity() *Entity {
	return NewEntityWithID(EntityID(uuid.New()))
}

func NewEntityWithID(id EntityID) *Entity {
	return &Entity{
		id:         id,
		components: make(map[ComponentType]Component),
	}
}

func (e *Entity) ID() EntityID {
	return e.id
}

func (e *Entity) AddComponents(components ...Component) (*Entity, error) {
	for _, component := range components {
		if component == nil {
			return nil, errors.NewNilError("cannot add nil component to entity")
		}
		if component.Type().IsNil() {
			return nil, errors.NewNilError("cannot add component with nil type to entity")
		}
		if _, exists := e.components[component.Type()]; exists {
			return nil, errors.NewDuplicateError("component already exists: " + component.Type().String())
		}
		e.components[component.Type()] = component
	}
	return e, nil
}

func (e *Entity) RemoveComponents(types ...ComponentType) (*Entity, error) {
	for _, t := range types {
		if t.IsNil() {
			return nil, errors.NewNilError("cannot remove component with nil type from entity")
		}
		if _, exists := e.components[t]; !exists {
			continue
		}
		delete(e.components, t)
	}
	return e, nil
}

func (e *Entity) GetComponent(t ComponentType) (Component, bool, error) {
	if t.IsNil() {
		return nil, false, errors.NewNilError("cannot get component with nil type from entity")
	}
	component, exists := e.components[t]
	if !exists {
		return nil, false, errors.NewNotFoundError("component not found: " + t.String())
	}
	return component, true, nil
}

func (e *Entity) HasComponents(types ...ComponentType) (bool, error) {
	for _, t := range types {
		if t.IsNil() {
			return false, errors.NewNilError("cannot check for component with nil type in entity")
		}
		if _, exists := e.components[t]; !exists {
			return false, nil
		}
	}
	return true, nil
}

func (e *Entity) AddTags(tags ...string) (*Entity, error) {
	for _, tag := range tags {
		if tag == "" {
			return nil, errors.NewNilError("cannot add empty tag to entity")
		}
		if e.tags.Contains(tag) {
			return nil, errors.NewDuplicateError("tag already exists: " + tag)
		}
		e.tags.Add(tag)
	}
	return e, nil
}

func (e *Entity) RemoveTags(tags ...string) (*Entity, error) {
	for _, tag := range tags {
		if tag == "" {
			return nil, errors.NewNilError("cannot remove empty tag from entity")
		}
		e.tags.Remove(tag)
	}
	return e, nil
}

func (e *Entity) HasTags(tags ...string) (bool, error) {
	for _, tag := range tags {
		if tag == "" {
			return false, errors.NewNilError("cannot check for empty tag in entity")
		}
		if !e.tags.Contains(tag) {
			return false, nil
		}
	}
	return true, nil
}
