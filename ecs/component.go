package ecs

import (
	"strconv"

	"github.com/adm87/finch-core/hash"
)

// ComponentType is a unique identifier for a component type.
//
// Used to identify components in the ECS framework.
type ComponentType hash.Hash

func (t ComponentType) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func (t ComponentType) IsNil() bool {
	return t == 0
}

type Component interface {
	// Type returns the unique identifier for the component type.
	Type() ComponentType
}
