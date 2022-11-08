package valdy

import (
	"fmt"
)

type (
	// CreateRuleFunc is supposed to be a wrapper func around NewRule func
	// in order to pass additional params to rule.
	// See "_example" folder for details.
	CreateRuleFunc[T any] func(T) Rule[T]
	// CheckFunc is function that validates passed target.
	CheckFunc[T any] func(target T) bool
)

// Rule is something that can be checked.
type Rule[T any] interface {
	Check(T) error
}

type rule[T any] struct {
	// Message represents the description of this rule.
	// If rule check is failed, then ErrValidation with the Message is returned.
	Message string
	check   CheckFunc[T]
}

// NewRule return new rule with passed message and check function.
func NewRule[T any](message string, check CheckFunc[T]) Rule[T] {
	return &rule[T]{
		Message: message,
		check:   check,
	}
}

// Check target against stored v.check function.
//
// If check is failed, then error with check's message will be returned.
func (v rule[T]) Check(target T) error {
	if ok := v.check(target); !ok {
		return v.buildError()
	}

	return nil
}

func (v rule[T]) buildError() error {
	if v.Message == "" {
		return fmt.Errorf("%w", ErrValidation)
	} else {
		return fmt.Errorf("%w: %s", ErrValidation, v.Message)
	}
}
