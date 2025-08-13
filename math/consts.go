package math

const (
	DegToRad = 0.017453292519943295 // Pi / 180
	RadToDeg = 57.29577951308232    // 180 / Pi
)

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
