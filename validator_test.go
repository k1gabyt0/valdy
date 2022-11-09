package valdy_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/k1gabyt0/valdy"
)

// personFixture is a struct for testing purposes.
type personFixture struct {
	name string
	age  int
}

func TestValidator_IncorrectMode(t *testing.T) {
	const INCORRECT_MODE = 999

	var validator valdy.Validator[string]
	validator.Mode = INCORRECT_MODE

	err := validator.Validate("something")
	if err == nil {
		t.Error("should be error if incorrect mode passed to validator")
	}
	if !errors.Is(err, valdy.ErrInternal) {
		t.Errorf("should be ErrInternal, but got=%v", err)
	}
}

func TestValidator_RunAllIsDefaultMode(t *testing.T) {
	var validator valdy.Validator[string]
	if validator.Mode != valdy.RUN_ALL {
		t.Errorf("expected RUN_ALL(%v) mode to be default, but default is %v", valdy.RUN_ALL, validator.Mode)
	}
}

func TestValidator_Validate_RunAllMode(t *testing.T) {
	const MODE = valdy.RUN_ALL

	type args struct {
		target personFixture
		rules  []valdy.CreateRuleFunc[personFixture]
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Validation against no rules always passes",
			args: args{
				target: personFixture{
					name: "John",
					age:  18,
				},
				rules: []valdy.CreateRuleFunc[personFixture]{},
			},
		},
		{
			name: "If all rules passed, then no error returned",
			args: args{
				target: personFixture{
					name: "John",
					age:  18,
				},
				rules: func() []valdy.CreateRuleFunc[personFixture] {
					var isAdult valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p personFixture) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p personFixture) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[personFixture]{isAdult, isNamedJohn}
				}(),
			},
		},
		{
			name: "If one rule failed, then validation error is returned",
			args: args{
				target: personFixture{
					name: "John",
					age:  17,
				},
				rules: func() []valdy.CreateRuleFunc[personFixture] {
					var isAdult valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p personFixture) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p personFixture) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[personFixture]{isAdult, isNamedJohn}
				}(),
			},
			wantErr: true,
		},
		{
			name: "If all rules failed, then validation error is returned",
			args: args{
				target: personFixture{
					name: "Not John =(",
					age:  17,
				},
				rules: func() []valdy.CreateRuleFunc[personFixture] {
					var isAdult = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p personFixture) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p personFixture) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[personFixture]{isAdult, isNamedJohn}
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var validator valdy.Validator[personFixture]
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

func TestValidator_Validate_StopOnFirstFailMode(t *testing.T) {
	const MODE = valdy.STOP_ON_FIRST_FAIL

	type args struct {
		target personFixture
		rules  []valdy.CreateRuleFunc[personFixture]
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Validation against no rules always passes",
			args: args{
				target: personFixture{
					name: "John",
					age:  18,
				},
				rules: []valdy.CreateRuleFunc[personFixture]{},
			},
		},
		{
			name: "If all rules passed, then no error returned",
			args: args{
				target: personFixture{
					name: "John",
					age:  18,
				},
				rules: func() []valdy.CreateRuleFunc[personFixture] {
					var isAdult valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p personFixture) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p personFixture) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[personFixture]{isAdult, isNamedJohn}
				}(),
			},
		},
		{
			name: "If one rule failed, then validation error is returned",
			args: args{
				target: personFixture{
					name: "John",
					age:  17,
				},
				rules: func() []valdy.CreateRuleFunc[personFixture] {
					var isAdult valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p personFixture) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn valdy.CreateRuleFunc[personFixture] = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p personFixture) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[personFixture]{isAdult, isNamedJohn}
				}(),
			},
			wantErr: true,
		},
		{
			name: "If all rules failed, then validation error is returned",
			args: args{
				target: personFixture{
					name: "Not John =(",
					age:  17,
				},
				rules: func() []valdy.CreateRuleFunc[personFixture] {
					var isAdult = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should be adult(18+)", target.name),
							func(p personFixture) bool {
								return p.age >= 18
							},
						)
					}
					var isNamedJohn = func(target personFixture) valdy.Rule[personFixture] {
						return valdy.NewRule(
							fmt.Sprintf("%s should has name 'John'", target.name),
							func(p personFixture) bool {
								return p.name == "John"
							},
						)
					}

					return []valdy.CreateRuleFunc[personFixture]{isAdult, isNamedJohn}
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var validator valdy.Validator[personFixture]
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
