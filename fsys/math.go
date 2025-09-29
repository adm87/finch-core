package fsys

const (
	HalfPi   = 1.5707963267948966   // π/2
	Pi       = 3.141592653589793    // π
	TwoPi    = 6.283185307179586    // 2π
	RadToDeg = 57.29577951308232    // 180/π
	DegToRad = 0.017453292519943295 // π/180
)

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func Abs[T number](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

func Clamp[T number](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func Lerp[T number](a, b, t T) T {
	return a + (b-a)*t
}

func Max[T number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}
