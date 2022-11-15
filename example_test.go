package valdy_test

import (
	"errors"
	"fmt"

	"github.com/k1gabyt0/valdy"
)

func Example() {
	var wantedCriminals = []string{
		"Ivan",
		"Jackob",
		// "John",
	}

	type person struct {
		name     string
		age      int
		children []person
	}

	// Validation rules:
	var (
		// ErrIsAdult is corresponding error for isAdult rule.
		// When isAdult check fails, the error(wrapped) is returned.
		// Thus, it is possible to check what rules have failed during validation process.
		ErrIsAdult = errors.New("person is not adult")
		isAdult    = func(p person) valdy.Checker[person] {
			return valdy.NewRule(
				fmt.Sprintf("%s must be adult. age is %d", p.name, p.age),
				func(p person) bool {
					return p.age >= 18
				},
			).WithError(ErrIsAdult)
		}

		ErrHasChildren = errors.New("person doesn't have children")
		hasChildren    = func(p person) valdy.Checker[person] {
			return valdy.NewRule(
				fmt.Sprintf("%s should has children, but he doesn't", p.name),
				func(p person) bool {
					return len(p.children) != 0
				},
			).WithError(ErrHasChildren)
		}

		ErrInWantedList = errors.New("person is not a wanted criminal")
		inWantedList    = func(criminalList []string) valdy.CreateRuleFunc[person] {
			return func(p person) valdy.Checker[person] {
				return valdy.NewRule(
					fmt.Sprintf("%s should be wanted criminal, but he is not in the list: %v", p.name, criminalList),
					func(p person) bool {
						for _, criminal := range criminalList {
							if p.name == criminal {
								return true
							}
						}
						return false
					},
				).WithError(ErrInWantedList)
			}
		}
	)

	var validator valdy.Validator[person]

	john := person{
		name: "John",
		age:  18,
	}

	err := validator.Validate(john,
		isAdult,
		hasChildren,
		inWantedList(wantedCriminals),
	)
	fmt.Println(err)

	// Also you can check which ones validations have failed.
	// In this case the thing you need to do is to proivde
	// corressponding error for each validation.
	// Otherwise, only the ErrValidation error can be checked,
	// but in this case the error has no idea what validation has failed.
	if errors.Is(err, ErrIsAdult) {
		fmt.Println("ErrIsAdult")
	}
	if errors.Is(err, ErrHasChildren) {
		fmt.Println("ErrHasChildren")
	}
	if errors.Is(err, ErrInWantedList) {
		fmt.Println("ErrInWantedList")
	}

	// Output: validation has failed:
	//	person doesn't have children:
	//	John should has children, but he doesn't
	//	person is not a wanted criminal:
	//	John should be wanted criminal, but he is not in the list: [Ivan Jackob]
	// ErrHasChildren
	// ErrInWantedList
}
