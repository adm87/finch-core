package linq

import "github.com/adm87/finch-core/types"

// Distinct returns a new slice containing only unique elements from the input slice.
func Distinct[T comparable](items []T) []T {
	if len(items) == 0 {
		return nil
	}

	seen := types.HashSet[T]{}
	var result []T

	for _, item := range items {
		if !seen.Contains(item) {
			seen.Add(item)
			result = append(result, item)
		}
	}

	return result
}

// DistinctFunc returns a new slice containing unique elements based on a key function.
func DistinctFunc[T any, K comparable](items []T, keyFunc func(T) K) []T {
	if len(items) == 0 {
		return nil
	}

	seen := types.HashSet[K]{}
	var result []T

	for _, item := range items {
		key := keyFunc(item)
		if !seen.Contains(key) {
			seen.Add(key)
			result = append(result, item)
		}
	}

	return result
}
