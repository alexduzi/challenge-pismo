package exception

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput                  = errors.New("invalid input")
	ErrDatabaseError                 = errors.New("database error")
	ErrInternalServerError           = errors.New("internal server error")
	ErrDuplicateDocument             = errors.New("document number already exists")
	ErrDuplicateEmail                = errors.New("email already exists")
	ErrNotFound                      = errors.New("resource not found")
	ErrInvalidAmountForOperationType = errors.New("invalid amount for operation type")
	ErrAmountGreaterThanZero         = errors.New("amount must be greater than 0")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrInvalidInput
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

func NewValidationErrorf(format string, args ...any) error {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}
