package hash

import (
	"hash/fnv"
	"reflect"
)

type Hash uint64

// GetHashFromType computes a hash from the type of T.
//
// Uses the type's package path and name to ensure uniqueness.
func GetHashFromType[T any]() Hash {
	t := reflect.TypeOf((*T)(nil)).Elem()

	var typeName string
	if t.PkgPath() != "" {
		typeName = t.PkgPath() + "." + t.Name()
	} else {
		typeName = t.String()
	}

	h := fnv.New64a()
	h.Write([]byte(typeName))

	return Hash(h.Sum64())
}
