package math

// Clamp restricts a value to be within a specified range.
func Clamp[T number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Sign returns -1 for negative values, 1 for positive values, and 0 for zero.
func Sign[T number](value T) int {
	if value < 0 {
		return -1
	} else if value > 0 {
		return 1
	}
	return 0
}
