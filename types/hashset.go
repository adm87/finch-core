package types

import "slices"

var entry = struct{}{}

type HashSet[T comparable] map[T]struct{}

func NewHashSetFromSlice[T comparable](items []T) HashSet[T] {
	hs := make(HashSet[T], len(items))
	for _, item := range items {
		hs[item] = entry
	}
	return hs
}

func (hs HashSet[T]) Add(item T) {
	hs[item] = entry
}

func (hs HashSet[T]) Remove(item T) {
	delete(hs, item)
}

func (hs HashSet[T]) Contains(item T) bool {
	_, exists := hs[item]
	return exists
}

func (hs HashSet[T]) Size() int {
	return len(hs)
}

func (hs HashSet[T]) Clear() {
	for k := range hs {
		delete(hs, k)
	}
}

func (hs HashSet[T]) ToSlice() []T {
	slice := make([]T, 0, len(hs))
	for k := range hs {
		slice = append(slice, k)
	}
	return slice
}

func (hs HashSet[T]) IsEmpty() bool {
	return len(hs) == 0
}
func Union[T comparable](sets ...HashSet[T]) HashSet[T] {
	if len(sets) == 0 {
		return make(HashSet[T])
	}
	if len(sets) == 1 {
		return sets[0]
	}

	result := make(HashSet[T])
	for _, set := range sets {
		for k := range set {
			result.Add(k)
		}
	}
	return result
}

func Intersection[T comparable](sets ...HashSet[T]) HashSet[T] {
	if len(sets) == 0 {
		return make(HashSet[T])
	}
	if len(sets) == 1 {
		return sets[0]
	}

	slices.SortFunc(sets, func(a, b HashSet[T]) int {
		return len(a) - len(b)
	})

	result := make(HashSet[T])
	for k := range sets[0] {
		inAll := true
		for _, set := range sets[1:] {
			if !set.Contains(k) {
				inAll = false
				break
			}
		}
		if inAll {
			result.Add(k)
		}
	}
	return result
}
