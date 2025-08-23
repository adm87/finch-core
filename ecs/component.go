package ecs

import (
	"strconv"

	"github.com/adm87/finch-core/utils"
)

// ComponentType is a unique identifier for a component type.
//
// Used to identify components in the ECS framework.
type ComponentType uint64

func (t ComponentType) IsNil() bool {
	return t == 0
}

func (t ComponentType) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func NewComponentType[T Component]() ComponentType {
	return ComponentType(utils.GetHashFromType[T]())
}

type Component interface {
	// Type returns the unique identifier for the component type.
	Type() ComponentType
}

func GetComponent[T Component](world *World, entity Entity, componentType ComponentType) (T, bool, error) {
	var zero T

	component, exists, err := world.GetComponent(entity, componentType)
	if err != nil {
		return zero, false, err
	}
	if !exists {
		return zero, false, nil
	}

	typedComponent, ok := component.(T)
	if !ok {
		return zero, false, ErrComponentTypeMismatch
	}

	return typedComponent, ok, nil
}
