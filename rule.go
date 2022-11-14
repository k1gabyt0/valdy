package valdy

type (
	// CreateRuleFunc is supposed to be a wrapper func around NewRule func
	// in order to pass checkTarget to rule.
	//
	// This approach allows to dynamically pass checkTarget's info into rule message.
	//
	// See "_example" folder for further details.
	CreateRuleFunc[T any] func(checkTarget T) Checker[T]
	// CheckFunc is a function that validates passed target.
	CheckFunc[T any] func(target T) bool
)

// Checker is something that can check and return appropriate error.
type Checker[T any] interface {
	Check(T) error
}

// Rule is an implementation of Checker.
//
// It can check target against passed check func.
//
// If check is failed it returns appropriate error.
// "Appropriate error" means that you can pass err to be wrapped into returned error.
// Returned error contains passed rule's Message.
// Function calls of errors.Is(appErr, ErrValidation) and errors.Is(appErr, err) will return true.
// If no err is passed, then only errors.Is(appErr, ErrValidation) will return true.
type Rule[T any] struct {
	// Message represents the description of this rule.
	// If rule check is failed, then ErrValidation with the Message is returned.
	Message string
	check   CheckFunc[T]
	err     error
}

// Make sure that Rule implements Checker.
var _ Checker[any] = &Rule[any]{}

// NewRule return new rule with passed message and check function.
func NewRule[T any](message string, check CheckFunc[T]) *Rule[T] {
	return &Rule[T]{
		Message: message,
		check:   check,
	}
}

// Check target against stored v.check function.
//
// If Check is failed, then error with check's message will be returned.
func (r *Rule[T]) Check(target T) error {
	if ok := r.check(target); !ok {
		valErr := NewValidationError(r.Message)
		if r.err == nil {
			return valErr
		}
		return ValidationErrorFrom(r.err, valErr)
	}

	return nil
}

// GetCheckFunc that holds rule validation logic.
func (r *Rule[T]) GetCheckFunc() CheckFunc[T] {
	return r.check
}

// WithError sets the err to be wrapped into resulting error when rule's Check is failed.
func (r *Rule[T]) WithError(err error) Checker[T] {
	r.err = err
	return r
}
