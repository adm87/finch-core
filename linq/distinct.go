package linq

import "github.com/adm87/finch-core/hashset"

// Distinct returns a new slice containing only unique elements from the input slice.
func Distinct[T comparable](items []T) []T {
	if len(items) == 0 {
		return nil
	}

	seen := hashset.New[T]()
	var result []T

	for _, item := range items {
		if !seen.Contains(item) {
			result = append(result, item)
		}
		seen.AddDistinct(item)
	}

	return result
}

// DistinctFunc returns a new slice containing unique elements based on a key function.
func DistinctFunc[T any, K comparable](items []T, keyFunc func(T) K) []T {
	if len(items) == 0 {
		return nil
	}

	seen := hashset.New[K]()
	var result []T

	for _, item := range items {
		key := keyFunc(item)
		if !seen.Contains(key) {
			result = append(result, item)
		}
		seen.AddDistinct(key)
	}

	return result
}
