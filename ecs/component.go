package ecs

import (
	"strconv"

	"github.com/adm87/finch-core/hash"
)

// ComponentType is a unique identifier for a component type.
//
// Used to identify components in the ECS framework.
type ComponentType hash.Hash

func (t ComponentType) IsNil() bool {
	return t == 0
}

func (t ComponentType) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func NewComponentType[T Component]() ComponentType {
	return ComponentType(hash.GetHashFromType[T]())
}

type Component interface {
	// Dispose can be used to dereference any resources held by the component.
	Dispose()

	// Type returns the unique identifier for the component type.
	Type() ComponentType
}
