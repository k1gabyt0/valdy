package valdy_test

import (
	"errors"
	"testing"

	"github.com/k1gabyt0/valdy"
)

func TestValidationError_Is(t *testing.T) {
	type args struct {
		msg      string
		original error
		errs     []error
	}

	errA := errors.New("error A")
	errB := errors.New("error B")

	tests := []struct {
		name          string
		args          args
		wantErrIs     []error
		dontWantErrIs []error
	}{
		{
			name: "Empty ValidationError is ErrValidation",
			args: args{
				msg:  "No errors inside",
				errs: []error{},
			},
			wantErrIs:     []error{valdy.ErrValidation},
			dontWantErrIs: []error{errA, errB},
		},
		{
			name: "ValidationError is unwrappable",
			args: args{
				msg:  "I have some errs",
				errs: []error{errA, errB},
			},
			wantErrIs:     []error{valdy.ErrValidation, errA, errB},
			dontWantErrIs: []error{},
		},
		{
			name: "Err that was 'transformed' into ValidationError is unwrappable",
			args: args{
				original: errA,
				errs:     []error{errB},
			},
			wantErrIs:     []error{valdy.ErrValidation, errA, errB},
			dontWantErrIs: []error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.args.original == nil {
				err = valdy.NewValidationError(tt.args.msg, tt.args.errs...)
			} else {
				err = valdy.ValidationErrorFrom(tt.args.original, tt.args.errs...)
			}

			if err == nil {
				t.Errorf("There is no ValidationError created")
				return
			}

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
