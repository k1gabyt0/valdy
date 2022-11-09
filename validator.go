package valdy

import (
	"fmt"
)

type ValidatorMode uint

const (
	// RUN_ALL means that all validations will be runned.
	// Even if some of them may be failed.
	RUN_ALL ValidatorMode = iota
	// STOP_ON_FIRST_FAIL means as soon as first validation fails,
	// others won't be checked.
	STOP_ON_FIRST_FAIL
)

// Validator validates anything you pass to it.
//
// It has different modes: RUN_ALL, STOP_ON_FIRST_FAIL.
// Zero-value of Validator uses RUN_ALL mode by default.
type Validator[T any] struct {
	// Mode that Validator uses for validation.
	//
	// The new ones can be added by embedding Validator struct and
	// overriding Validate func.
	Mode ValidatorMode
}

// Validate checks target against passed validation rules creators.
//
// If Validator has incorrect Mode, then Validate will return ErrInternal.
func (v *Validator[T]) Validate(target T, creators ...CreateRuleFunc[T]) error {
	rules := extractRules(target, creators)

	switch v.Mode {
	case RUN_ALL:
		return v.runAllValidation(target, rules)
	case STOP_ON_FIRST_FAIL:
		return v.stopOnFirstFailValidation(target, rules)
	default:
		return fmt.Errorf("%w: no such mode %d for validator", ErrInternal, v.Mode)
	}
}

func (v *Validator[T]) runAllValidation(target T, rules []Rule[T]) error {
	var lastFailedErr error
	for _, rule := range rules {
		if err := rule.Check(target); err != nil {
			lastFailedErr = err
		}
	}
	return lastFailedErr
}

func (v *Validator[T]) stopOnFirstFailValidation(target T, rules []Rule[T]) error {
	for _, rule := range rules {
		if err := rule.Check(target); err != nil {
			return err
		}
	}
	return nil
}

// extractRules extracts validation rules from creators.
func extractRules[T any](validationTarget T, creators []CreateRuleFunc[T]) []Rule[T] {
	rules := make([]Rule[T], 0, len(creators))
	for _, c := range creators {
		rules = append(rules, c(validationTarget))
	}
	return rules
}
