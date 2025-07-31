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

func (e *Entity) AddComponents(components ...Component) *Entity {
	for _, component := range components {
		if component == nil {
			panic(errors.NewNilError("cannot add nil component to entity"))
		}
		if component.Type().IsNil() {
			panic(errors.NewNilError("cannot add component with nil type to entity"))
		}
		if _, exists := e.components[component.Type()]; exists {
			panic(errors.NewDuplicateError("component already exists: " + component.Type().String()))
		}
		e.components[component.Type()] = component
	}
	return e
}

func (e *Entity) RemoveComponents(types ...ComponentType) *Entity {
	for _, t := range types {
		if t.IsNil() {
			panic(errors.NewNilError("cannot remove component with nil type from entity"))
		}
		if _, exists := e.components[t]; !exists {
			continue
		}
		delete(e.components, t)
	}
	return e
}

func (e *Entity) HasComponents(types ...ComponentType) bool {
	for _, t := range types {
		if t.IsNil() {
			panic(errors.NewNilError("cannot check for component with nil type in entity"))
		}
		if _, exists := e.components[t]; !exists {
			return false
		}
	}
	return true
}

func (e *Entity) AddTags(tags ...string) *Entity {
	for _, tag := range tags {
		if tag == "" {
			panic(errors.NewNilError("cannot add empty tag to entity"))
		}
		e.tags.Add(tag)
	}
	return e
}

func (e *Entity) RemoveTags(tags ...string) *Entity {
	for _, tag := range tags {
		if tag == "" {
			panic(errors.NewNilError("cannot remove empty tag from entity"))
		}
		e.tags.Remove(tag)
	}
	return e
}

func (e *Entity) HasTags(tags ...string) bool {
	for _, tag := range tags {
		if tag == "" {
			panic(errors.NewNilError("cannot check for empty tag in entity"))
		}
		if !e.tags.Contains(tag) {
			return false
		}
	}
	return true
}
