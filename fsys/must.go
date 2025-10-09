package fsys

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustGet[T any](val T, err error) T {
	Must(err)
	return val
}
