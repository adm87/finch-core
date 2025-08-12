package ecs

import (
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

var NilEntity = Entity(uuid.Nil)
