package valdy_test

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/k1gabyt0/valdy"
)

func TestRule_Check(t *testing.T) {
	type args struct {
		message string
		check   valdy.CheckFunc[personFixture]
		withErr error
	}

	wrongName := errors.New("wrong name error")

	tests := []struct {
		name      string
		toCheck   personFixture
		args      args
		wantErr   bool
		wantErrIs error
	}{
		{
			name: "If rule check with concrete error fails, then return rule corresponding error",
			toCheck: personFixture{
				name: "John",
			},
			args: args{
				message: "Name should not be John!",
				check: func(p personFixture) bool {
					return p.name != "John"
				},
				withErr: wrongName,
			},
			wantErr:   true,
			wantErrIs: wrongName,
		},
		{
			name: "If rule check fails, then return rule error",
			toCheck: personFixture{
				name: "John",
			},
			args: args{
				message: "Name should not be John!",
				check: func(p personFixture) bool {
					return p.name != "John"
				},
			},
			wantErr: true,
		},
		{
			name: "Empty rule message is fine too",
			toCheck: personFixture{
				name: "John",
			},
			args: args{
				message: "",
				check: func(p personFixture) bool {
					return p.name != "John"
				},
			},
			wantErr: true,
		},
		{
			name: "If rule check passes, then no error",
			toCheck: personFixture{
				name: "Mikhail",
			},
			args: args{
				message: "Name should be >=5 letters long",
				check: func(target personFixture) bool {
					return utf8.RuneCountInString(target.name) >= 5
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := valdy.
				NewRule(tt.args.message, tt.args.check).
				WithError(tt.args.withErr)

			err := rule.Check(tt.toCheck)
			if err == nil {
				if tt.wantErr || tt.wantErrIs != nil {
					t.Error("wanted error, but didn't get it")
				}
			}

			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error=%q, but didn't wanted it", err)
				}
				if tt.wantErrIs != nil {
					if !errors.Is(err, tt.wantErrIs) {
						t.Errorf("got error=%q, but it is not %q", err, tt.wantErrIs)
					}
				}
				if !errors.Is(err, valdy.ErrValidation) {
					t.Errorf("got error=%q, but it is not ErrValidation", err)
				}

				if !strings.Contains(err.Error(), tt.args.message) {
					t.Errorf("got error=%q, but it doesn't contain passed description %q", err, tt.args.message)
				}
			}
		})
	}
}
