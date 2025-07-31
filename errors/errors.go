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
