package valdy

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrValidation shows that validation is failed.
	ErrValidation = errors.New("validation error")
	// ErrInternal shows that non-validation error happened.
	ErrInternal = errors.New("internal error")
)

// ValidationError is error that able to store multiple errors.
//
// Whether ValidationError is checked by errors.Is(err, ErrValidation), it will be true.
type ValidationError struct {
	original error
	msg      string
	errs     []error
}

// Make sure that ValidationError is an error.
var _ error = &ValidationError{}

// From creates a ValidationError err from passed original error
// and it's inner errs.
//
// err's msg will be the same as original.Error().
// Function call errors.Is(err, original) will return true.
//
// If passed original is already a ValidationError, then original is returned.
func From(original error, errs ...error) *ValidationError {
	if original == nil {
		return NewValidationError("")
	}

	if tgt, ok := original.(*ValidationError); ok {
		return tgt
	}

	e := &ValidationError{
		original: original,
		msg:      original.Error(),
	}
	return e.WithErrors(errs...)
}

// NewValidationError returns ValidationError with passed msg and errs.
//
// Passed nil errs though will be filtered. Only non-nil errs are stored.
func NewValidationError(msg string, errs ...error) *ValidationError {
	e := &ValidationError{
		msg: msg,
	}
	return e.WithErrors(errs...)
}

// WithErrors stores passed errs in ValidationError.
//
// Passed nil errs though will be filtered. Only non-nil errs are stored.
func (v *ValidationError) WithErrors(errs ...error) *ValidationError {
	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}
	v.errs = nonNilErrs
	return v
}

// Error formats error messages with new lines starting with v.msg.
func (v *ValidationError) Error() string {
	var msgBuilder strings.Builder
	msgBuilder.WriteString(v.msg)

	if len(v.errs) > 0 {
		msgBuilder.WriteString(":")

		for _, err := range v.errs {
			msgBuilder.WriteString(fmt.Sprintf("\n\t%s", err.Error()))
		}
	}

	return msgBuilder.String()
}

// Is reports whether any error in v's tree matches target.
//
// It is always returns true if checked against ErrValidation.
func (v *ValidationError) Is(target error) bool {
	if target == ErrValidation || target == v.original {
		return true
	}

	if tgt, ok := target.(*ValidationError); ok {
		isSameMessage := v.msg == tgt.msg
		hasSameErrors := equalErrors(v.errs, tgt.errs)

		if isSameMessage && hasSameErrors {
			return true
		}
	}

	for _, err := range v.errs {
		if is := errors.Is(err, target); is {
			return true
		}
	}
	return false
}

// As finds the first error in v's tree that matches target(starting with the root),
// and if one is found, sets target to that error value and returns true.
// Otherwise, it returns false.
//
// As finds the first matching error in a preorder traversal of the tree.
func (v *ValidationError) As(target any) bool {
	if errors.As(v.original, target) {
		return true
	}

	for _, err := range v.errs {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

// Message returns root message of ValdiationError that will be displayed first.
func (v *ValidationError) Message() string {
	return v.msg
}

// Original error that have been passed to become ValidationError.
func (v *ValidationError) Original() error {
	return v.original
}

// Errors returns all errors stored in v.
func (v *ValidationError) Errors() []error {
	return v.errs
}

func equalErrors(s1, s2 []error) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
