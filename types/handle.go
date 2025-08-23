package types

type Handle[T any] interface {
	Get() (T, error)
	IsValid() bool
}
