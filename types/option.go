package types

type Optional[T any] struct {
	value T
	valid bool
}

func NewOption[T any](value T) Optional[T] {
	return Optional[T]{value: value, valid: true}
}

func NewEmptyOption[T any]() Optional[T] {
	return Optional[T]{valid: false}
}

func (o Optional[T]) IsValid() bool {
	return o.valid
}

func (o Optional[T]) Value() T {
	return o.value
}

func (o *Optional[T]) SetValue(value T) {
	o.value = value
	o.valid = true
}

func (o *Optional[T]) Invalidate() {
	var zero T

	o.value = zero
	o.valid = false
}
