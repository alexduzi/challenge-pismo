package exception

import "errors"

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrDatabaseError       = errors.New("database error")
	ErrInternalServerError = errors.New("internal server error")
)
