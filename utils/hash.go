package utils

import (
	"hash/fnv"
	"reflect"
)

func GetHashFromType[T any]() uint64 {
	t := reflect.TypeOf((*T)(nil)).Elem()

	var typeName string
	if t.PkgPath() != "" {
		typeName = t.PkgPath() + "." + t.Name()
	} else {
		typeName = t.String()
	}

	h := fnv.New64a()
	h.Write([]byte(typeName))

	return h.Sum64()
}
