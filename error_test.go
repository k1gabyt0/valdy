package valdy_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/k1gabyt0/valdy"
)

func TestValidationError_Equality(t *testing.T) {
	err1 := valdy.NewValidationError("error epta")
	err2 := valdy.NewValidationError("error epta")
	if err1 == err2 {
		t.Errorf("errors(%q, %q) with same message should not be equal", err1, err2)
	}
}

func TestValidationError_From(t *testing.T) {
	type args struct {
		original error
		children []error
	}

	ErrA := errors.New("error A")
	ErrB := errors.New("error B")
	ErrC := errors.New("error C")
	ErrValA := valdy.NewValidationError("validation error", ErrA)

	tests := []struct {
		name      string
		args      args
		wantMsg   string
		wantIs    []error
		wantIsNot []error
	}{
		{
			name: "If original is nil, then return empty ValidationError",
			args: args{
				original: nil,
			},
			wantMsg: "",
			wantIs:  []error{valdy.ErrValidation},
		},
		{
			name: "Passed error is simple",
			args: args{
				original: ErrA,
			},
			wantMsg: ErrA.Error(),
			wantIs:  []error{valdy.ErrValidation, ErrA},
		},
		{
			name: "Passed error and children",
			args: args{
				original: ErrA,
				children: []error{ErrB, ErrC},
			},
			wantMsg: fmt.Sprintf("%s:\n\t%s\n\t%s", ErrA.Error(), ErrB.Error(), ErrC),
			wantIs:  []error{valdy.ErrValidation, ErrA, ErrB, ErrC},
		},
		{
			name: "Passed error and children(some of them nil)",
			args: args{
				original: ErrA,
				children: []error{
					nil,
					ErrB,
					nil,
					ErrC,
					nil,
					nil,
				},
			},
			wantMsg: fmt.Sprintf("%s:\n\t%s\n\t%s", ErrA.Error(), ErrB.Error(), ErrC),
			wantIs:  []error{valdy.ErrValidation, ErrA, ErrB, ErrC},
		},
		{
			name: "Passed error is another ValidationError",
			args: args{
				original: ErrValA,
			},
			wantMsg: ErrValA.Error(),
			wantIs:  []error{valdy.ErrValidation, ErrValA, ErrA},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valdy.From(tt.args.original, tt.args.children...)
			if err == nil {
				t.Error("created ValidationError is nil")
				return
			}

			if tt.wantMsg != err.Error() {
				t.Errorf("wanted message to be=%q, but got=%q", tt.wantMsg, err.Error())
			}

			for _, wantErr := range tt.wantIs {
				if !errors.Is(err, wantErr) {
					t.Errorf("expected that %q is %q", err, wantErr)
				}
			}
			for _, dontWantErr := range tt.wantIsNot {
				if errors.Is(err, dontWantErr) {
					t.Errorf("expected that %q is NOT %q", err, dontWantErr)
				}
			}
		})
	}
}

func TestValidationError_NewValidationErrorFiltersNilErrs(t *testing.T) {
	type args struct {
		errs []error
	}

	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")

	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{
			name: "No errs passed - no errs stored",
			args: args{
				errs: []error{},
			},
			wantErrs: []error{},
		},
		{
			name: "All nil errs passed - no errs stored",
			args: args{
				errs: []error{nil, nil, nil},
			},
			wantErrs: []error{},
		},
		{
			name: "Non nil errs passed - non nil errs stored",
			args: args{
				errs: []error{
					err1,
					err2,
					err3,
				},
			},
			wantErrs: []error{
				err1,
				err2,
				err3,
			},
		},
		{
			name: "Non nil & nil errs passed - non nil errs stored",
			args: args{
				errs: []error{
					err1,
					nil,
					err2,
					nil,
					err3,
					nil,
				},
			},
			wantErrs: []error{
				err1,
				err2,
				err3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valdy.NewValidationError("error msg", tt.args.errs...)
			if err == nil {
				t.Error("no ValidationError was created")
				return
			}

			errs := err.Errors()
			if len(errs) == len(tt.wantErrs) {
				for _, e := range errs {
					var found bool
					for _, wantE := range tt.wantErrs {
						if wantE == e {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("returned stored error=%q is not int wantErr=%v", e, tt.wantErrs)
					}
				}
			} else {
				t.Errorf("errs we want=%v and errs we got=%v have different sizes", tt.wantErrs, errs)
			}
		})
	}
}

func TestValidationError_MessageFormat(t *testing.T) {
	type args struct {
		msg       string
		innerErrs []error
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should only show passed err message if no inner errors",
			args: args{
				msg:       "Validation error with no inner errors",
				innerErrs: []error{},
			},
		},
		{
			name: "Should show passed err message and all messages for all inner errors ",
			args: args{
				msg: "Validation error with no inner errors",
				innerErrs: []error{
					errors.New("error 1"),
					errors.New("error 2"),
					errors.New("error 3"),
					errors.New("error 4"),
				},
			},
		},
		{
			name: "Empty message is fine too",
			args: args{
				msg: "",
				innerErrs: []error{
					errors.New("error 1"),
					errors.New("error 2"),
					errors.New("error 3"),
					errors.New("error 4"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valdy.NewValidationError(tt.args.msg, tt.args.innerErrs...)
			if err == nil {
				t.Errorf("there is no error created")
				return
			}
			if !strings.Contains(err.Error(), tt.args.msg) {
				t.Errorf("error doesn't contain passed message=%s", tt.args.msg)
			}
			for _, innerErr := range tt.args.innerErrs {
				if !strings.Contains(err.Error(), innerErr.Error()) {
					t.Errorf("resulting error message=%s doesn't contain inner error's=%q message=%q", err.Error(), innerErr, innerErr.Error())
				}
			}
		})
	}
}

func TestValidationError_Is(t *testing.T) {
	errA := errors.New("validation A failed")
	errB := errors.New("validation B failed")
	errC := errors.New("validation C failed")

	errComplexBAndC := valdy.NewValidationError("errComplexBAndC", errB, errC)
	errSuperComplex := valdy.NewValidationError("errSuperComplex", errComplexBAndC)

	type args struct {
		innerErrs []error
	}

	tests := []struct {
		name          string
		args          args
		wantErrIs     []error
		dontWantErrIs []error
	}{
		{
			name: "No inner errors",
			args: args{
				innerErrs: []error{},
			},
			wantErrIs:     []error{valdy.ErrValidation},
			dontWantErrIs: []error{errA, errB, errC},
		},
		{
			name: "One inner error",
			args: args{
				innerErrs: []error{errA},
			},
			wantErrIs:     []error{valdy.ErrValidation, errA},
			dontWantErrIs: []error{errB, errC},
		},
		{
			name: "Many inner errors",
			args: args{
				innerErrs: []error{errA, errB, errC},
			},
			wantErrIs:     []error{valdy.ErrValidation, errA, errB, errC},
			dontWantErrIs: []error{},
		},
		{
			name: "Complex inner error with some errors",
			args: args{
				innerErrs: []error{errSuperComplex},
			},
			wantErrIs:     []error{valdy.ErrValidation, errB, errC, errComplexBAndC, errSuperComplex},
			dontWantErrIs: []error{errA},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valdy.NewValidationError("validation error", tt.args.innerErrs...)
			for _, e := range tt.wantErrIs {
				if !errors.Is(err, e) {
					t.Errorf("expected err=(%q) to be (%q)", err, e)
				}
			}
			for _, e := range tt.dontWantErrIs {
				if errors.Is(err, e) {
					t.Errorf("not expected err=(%q) to be (%q)", err, e)
				}
			}
		})
	}
}

type simpleValidationError struct {
	message string
}

func (v simpleValidationError) Error() string {
	return v.message
}

func TestValidationError_As(t *testing.T) {
	errB := &simpleValidationError{message: "validation B failed"}
	errC := &simpleValidationError{message: "validation C failed"}
	errComplexBAndC := valdy.NewValidationError("errComplexBAndC", errB, errC)

	type args struct {
		original  error
		innerErrs []error
	}

	tests := []struct {
		name          string
		args          args
		targetsOkFn   func() *simpleValidationError
		targetsFailFn func() *simpleValidationError
	}{
		{
			name: "Should not be setted to unrelated type",
			args: args{},
			targetsFailFn: func() *simpleValidationError {
				var unrelated *simpleValidationError
				return unrelated
			},
		},
		{
			name: "Should be setted to related type",
			args: args{
				innerErrs: []error{errB},
			},
			targetsOkFn: func() *simpleValidationError {
				var related *simpleValidationError
				return related
			},
		},
		{
			name: "Should be setted to related type in complex case",
			args: args{
				innerErrs: []error{errComplexBAndC},
			},
			targetsOkFn: func() *simpleValidationError {
				var related *simpleValidationError
				return related
			},
		},
		{
			name: "Should be setted to original type",
			args: args{
				original: errB,
			},
			targetsOkFn: func() *simpleValidationError {
				var original *simpleValidationError
				return original
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.args.original != nil {
				err = valdy.From(tt.args.original, tt.args.innerErrs...)
			} else {
				err = valdy.NewValidationError("validation error", tt.args.innerErrs...)
			}

			errType := reflect.TypeOf(err)

			// should be setted to itself
			var selfErr *valdy.ValidationError
			if !errors.As(err, &selfErr) {
				t.Errorf("expected %q to be setted to itself", errType)
			}

			if tt.targetsOkFn != nil {
				targetOk := tt.targetsOkFn()
				trgtType := reflect.TypeOf(targetOk)
				if !errors.As(err, &targetOk) {
					t.Errorf("expected %q to be setted into %q", errType, trgtType)
				}
			}

			if tt.targetsFailFn != nil {
				targetFail := tt.targetsFailFn()
				trgtType := reflect.TypeOf(targetFail)
				if errors.As(err, &targetFail) {
					t.Errorf("expected %q to be NOT setted into %q", errType, trgtType)
				}
			}
		})
	}
}
