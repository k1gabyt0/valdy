package valdy

import (
	"errors"
)

var (
	ErrValidation = errors.New("validation error")
	ErrInternal   = errors.New("internal validation error")
)
