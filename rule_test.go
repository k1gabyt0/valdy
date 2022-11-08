package valdy_test

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/k1gabyt0/valdy"
)

func TestRule_Check(t *testing.T) {
	// person is a testing fixture
	type person struct {
		name string
	}

	type args struct {
		message string
		check   valdy.CheckFunc[person]
	}

	tests := []struct {
		name    string
		toCheck person
		args    args
		wantErr bool
	}{
		{
			name: "If rule check fails, then return rule error",
			toCheck: person{
				name: "John",
			},
			args: args{
				message: "Name should not be John!",
				check: func(p person) bool {
					return p.name != "John"
				},
			},
			wantErr: true,
		},
		{
			name: "Empty rule message is fine too",
			toCheck: person{
				name: "John",
			},
			args: args{
				message: "",
				check: func(p person) bool {
					return p.name != "John"
				},
			},
			wantErr: true,
		},
		{
			name: "If rule check passes, then no error",
			toCheck: person{
				name: "Mikhail",
			},
			args: args{
				message: "Name should be >=5 letters long",
				check: func(target person) bool {
					return utf8.RuneCountInString(target.name) >= 5
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := valdy.NewRule(tt.args.message, tt.args.check)
			err := rule.Check(tt.toCheck)
			if err == nil && tt.wantErr {
				t.Error("wanted error, but didn't get it")
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error=%v, but didn't wanted it", err)
				}
				if !errors.Is(err, valdy.ErrValidation) {
					t.Errorf("got error=%v, but it is not ErrValidation", err)
				}
				if !strings.Contains(err.Error(), tt.args.message) {
					t.Errorf("got error=%v, but it doesn't contain passed description %s", err, tt.args.message)
				}
			}
		})
	}
}
