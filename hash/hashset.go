package hash

type HashSet[T comparable] map[T]struct{}

func (s HashSet[T]) Add(item T) {
	s[item] = struct{}{}
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
