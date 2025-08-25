package linq

func Fill[T any](arr []T, value T) {
	for i := 0; i < len(arr); i++ {
		arr[i] = value
	}
}

func FillFunc[T any](arr []T, f func(int) T) {
	for i := 0; i < len(arr); i++ {
		arr[i] = f(i)
	}
}
