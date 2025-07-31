package errors

// =================================================================
// Duplicate Error
// =================================================================

type DuplicateError string

func (e DuplicateError) Error() string {
	return string(e)
}

func IsDuplicateError(err error) bool {
	_, ok := err.(DuplicateError)
	return ok
}

func NewDuplicateError(msg string) DuplicateError {
	return DuplicateError("duplicate: " + msg)
}

// =================================================================
// Invalid Argument Error
// =================================================================

type InvalidArgumentError string

func (e InvalidArgumentError) Error() string {
	return string(e)
}

func IsInvalidArgumentError(err error) bool {
	_, ok := err.(InvalidArgumentError)
	return ok
}

func NewInvalidArgumentError(msg string) InvalidArgumentError {
	return InvalidArgumentError("invalid argument: " + msg)
}

// =================================================================
// Nil Error
// =================================================================

type NilError string

func (e NilError) Error() string {
	return string(e)
}

func IsNilError(err error) bool {
	_, ok := err.(NilError)
	return ok
}

func NewNilError(msg string) NilError {
	return NilError("nil: " + msg)
}
