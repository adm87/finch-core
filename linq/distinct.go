package linq

import "github.com/adm87/finch-core/hash"

// Distinct returns a new slice containing only unique elements from the input slice.
func Distinct[T comparable](items []T) []T {
	if len(items) == 0 {
		return nil
	}

	seen := hash.HashSet[T]{}
	var result []T

	for _, item := range items {
		if !seen.Contains(item) {
			seen.Add(item)
			result = append(result, item)
		}
	}

	return result
}

// DistinctBy returns a new slice containing unique elements based on a key function.
func DistinctBy[T any, K comparable](items []T, keyFunc func(T) K) []T {
	if len(items) == 0 {
		return nil
	}

	seen := hash.HashSet[K]{}
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
