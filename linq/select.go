package linq

func Select[T comparable](source []T, selector func(T) bool) []T {
	var result []T
	for _, item := range source {
		if selector(item) {
			result = append(result, item)
		}
	}
	return result
}

func SelectKeys[K comparable, V any](source map[K]V, selector func(K, V) bool) []K {
	var result []K
	for key, value := range source {
		if selector(key, value) {
			result = append(result, key)
		}
	}
	return result
}

func SelectValues[K comparable, V any](source map[K]V, selector func(K, V) bool) []V {
	var result []V
	for key, value := range source {
		if selector(key, value) {
			result = append(result, value)
		}
	}
	return result
}

func SelectKVP[K comparable, V any](source map[K]V, selector func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for key, value := range source {
		if selector(key, value) {
			result[key] = value
		}
	}
	return result
}
