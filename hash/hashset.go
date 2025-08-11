package hash

import "slices"

var SetEntry struct{} = struct{}{}

type HashSet[T comparable] map[T]struct{}

func (s HashSet[T]) Add(item T) {
	s[item] = SetEntry
}

func (s HashSet[T]) Contains(item T) bool {
	_, exists := s[item]
	return exists
}

func (s HashSet[T]) Remove(item T) {
	delete(s, item)
}

func (s HashSet[T]) IsEmpty() bool {
	return len(s) == 0
}

func (s HashSet[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func MakeSetFrom[T comparable](items ...T) HashSet[T] {
	set := make(HashSet[T])
	for _, item := range items {
		set.Add(item)
	}
	return set
}

func IntersectHashSets[T comparable](sets ...HashSet[T]) HashSet[T] {
	if len(sets) == 0 {
		return make(HashSet[T])
	}

	// We'll sort the sets by smallest to largest to minimize work
	slices.SortStableFunc(sets, func(a, b HashSet[T]) int {
		return len(a) - len(b)
	})

	intersectionSet := make(HashSet[T])
	for item := range sets[0] {
		exists := true
		for i := 1; i < len(sets); i++ {
			if !sets[i].Contains(item) {
				exists = false
				break
			}
		}
		if exists {
			intersectionSet.Add(item)
		}
	}
	return intersectionSet
}
