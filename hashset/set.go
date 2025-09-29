package hashset

type Entry struct{}

type Set[T comparable] map[T]Entry

func New[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(item ...T) {
	for _, i := range item {
		s[i] = Entry{}
	}
}

func (s Set[T]) AddDistinct(items ...T) {
	for _, item := range items {
		if !s.Contains(item) {
			s[item] = Entry{}
		}
	}
}

func (s Set[T]) Remove(item ...T) {
	for _, i := range item {
		delete(s, i)
	}
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for k := range s {
		result = append(result, k)
	}
	return result
}
