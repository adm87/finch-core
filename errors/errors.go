package errors

// =================================================================
// Ambiguous Error
// =================================================================

type AmbiguousError string

func (e AmbiguousError) Error() string {
	return string(e)
}

func IsAmbiguousError(err error) bool {
	_, ok := err.(AmbiguousError)
	return ok
}

func NewAmbiguousError(msg string) AmbiguousError {
	return AmbiguousError("ambiguous: " + msg)
}

// =================================================================
// Conflict Error
// =================================================================

type ConflictError string

func (e ConflictError) Error() string {
	return string(e)
}

func IsConflictError(err error) bool {
	_, ok := err.(ConflictError)
	return ok
}

func NewConflictError(msg string) ConflictError {
	return ConflictError("conflict: " + msg)
}

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
// IO Error
// =================================================================

type IOError string

func (e IOError) Error() string {
	return string(e)
}

func IsIOError(err error) bool {
	_, ok := err.(IOError)
	return ok
}

func NewIOError(msg string) IOError {
	return IOError("io: " + msg)
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
// Not Implemented Error
// =================================================================

type NotImplementedError string

func (e NotImplementedError) Error() string {
	return string(e)
}

func IsNotImplementedError(err error) bool {
	_, ok := err.(NotImplementedError)
	return ok
}

func NewNotImplementedError(msg string) NotImplementedError {
	return NotImplementedError("not implemented: " + msg)
}

// =================================================================
// Parallel Error
// =================================================================

type ParallelError string

func (e ParallelError) Error() string {
	return string(e)
}

func IsParallelError(err error) bool {
	_, ok := err.(ParallelError)
	return ok
}

func NewParallelError(msg string) ParallelError {
	return ParallelError("parallel: " + msg)
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
