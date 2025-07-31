package hash

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
