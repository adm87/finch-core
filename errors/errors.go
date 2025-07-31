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

// =================================================================
// Not Found Error
// =================================================================

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

func NewNotFoundError(msg string) NotFoundError {
	return NotFoundError("not found: " + msg)
}

// =================================================================
// Security Error
// =================================================================

type SecurityError string

func (e SecurityError) Error() string {
	return string(e)
}

func IsSecurityError(err error) bool {
	_, ok := err.(SecurityError)
	return ok
}

func NewSecurityError(msg string) SecurityError {
	return SecurityError("security: " + msg)
}
