package linq

// Batch splits a slice into batches of a specified size.
func Batch[T any](items []T, size int) [][]T {
	n := len(items)

	if n == 0 || size <= 0 {
		return nil
	}

	if size >= n {
		return [][]T{items}
	}

	batchCount := (n + size - 1) / size
	batchSize := n / batchCount
	remainder := n % batchCount

	var batches [][]T
	start := 0

	for i := 0; i < batchCount; i++ {
		count := batchSize
		if i < remainder {
			count++
		}
		end := start + count
		batches = append(batches, items[start:end])
		start = end
	}

	return batches
}
