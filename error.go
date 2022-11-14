package valdy

import (
	"errors"

	"github.com/k1gabyt0/erry"
)

var (
	// ErrValidation shows that validation is failed.
	ErrValidation = errors.New("validation error")
	// ErrInternal shows that non-validation error happened.
	ErrInternal = errors.New("internal error")
)

// ValidationError is an error that able to store multiple errors.
//
// Whether ValidationError is checked by errors.Is(err, ErrValidation), it will be true.
type ValidationError struct {
	*erry.MError
}

func NewValidationError(msg string, errs ...error) *ValidationError {
	return &ValidationError{
		MError: erry.NewError(msg, errs...),
	}
}

func ValidationErrorFrom(original error, errs ...error) *ValidationError {
	return &ValidationError{
		MError: erry.ErrorFrom(original, errs...),
	}
}

// Is reports whether any error in v's tree matches target.
//
// It is always returns true if checked against ErrValidation.
func (v ValidationError) Is(target error) bool {
	if target == ErrValidation {
		return true
	}
	return v.MError.Is(target)
}
