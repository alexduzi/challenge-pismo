package exception

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{Message: "test error message"}

	if err.Error() != "test error message" {
		t.Errorf("expected 'test error message', got '%s'", err.Error())
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	err := &ValidationError{Message: "test"}

	if !errors.Is(err, ErrInvalidInput) {
		t.Error("expected ValidationError to unwrap to ErrInvalidInput")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("validation failed")

	if err.Error() != "validation failed" {
		t.Errorf("expected 'validation failed', got '%s'", err.Error())
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Error("expected error to be of type *ValidationError")
	}
}

func TestNewValidationErrorf(t *testing.T) {
	err := NewValidationErrorf("field %s is invalid: %d", "amount", 42)

	expected := "field amount is invalid: 42"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Error("expected error to be of type *ValidationError")
	}
}
