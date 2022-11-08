package valdy_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/k1gabyt0/valdy"
)

func TestValidator_Validate_StopOnFirstFailMode(t *testing.T) {
	const MODE = valdy.STOP_ON_FIRST_FAIL

	// person is a testing fixture
	type person struct {
		name string
		age  int
	}

	type args struct {
		target person
		rules  []valdy.CreateRuleFunc[person]
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Validation against no rules always passes",
			args: args{
				target: person{
					name: "John",
					age:  18,
				},
				rules: []valdy.CreateRuleFunc[person]{},
			},
		},
		{
			name: "If all rules passed, then no error returned",
			args: args{
				target: person{
					name: "John",
					age:  18,
				},
				rules: func() []valdy.CreateRuleFunc[person] {
					var isAdult valdy.CreateRuleFunc[person] = func(target person) valdy.Rule[person] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p person) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn valdy.CreateRuleFunc[person] = func(target person) valdy.Rule[person] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p person) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[person]{isAdult, isNamedJohn}
				}(),
			},
		},
		{
			name: "If one rule failed, then validation error is returned",
			args: args{
				target: person{
					name: "John",
					age:  17,
				},
				rules: func() []valdy.CreateRuleFunc[person] {
					var isAdult valdy.CreateRuleFunc[person] = func(target person) valdy.Rule[person] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p person) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn valdy.CreateRuleFunc[person] = func(target person) valdy.Rule[person] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p person) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[person]{isAdult, isNamedJohn}
				}(),
			},
			wantErr: true,
		},
		{
			name: "If all rules failed, then validation error is returned",
			args: args{
				target: person{
					name: "Not John =(",
					age:  17,
				},
				rules: func() []valdy.CreateRuleFunc[person] {
					var isAdult = func(target person) valdy.Rule[person] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p person) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn = func(target person) valdy.Rule[person] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p person) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[person]{isAdult, isNamedJohn}
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var validator valdy.Validator[person]
			validator.Mode = MODE

			err := validator.Validate(tt.args.target, tt.args.rules...)
			if err == nil && tt.wantErr {
				t.Error("wanted error, but didn't get it")
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error=%v, but didn't asked for this", err)
				}
				if !errors.Is(err, valdy.ErrValidation) {
					t.Errorf("got error=%v, but wanted %v", err, valdy.ErrValidation)
				}
			}
		})
	}
}
