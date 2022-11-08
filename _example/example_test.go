package example_test

import (
	"fmt"

	"github.com/k1gabyt0/valdy"
)

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
	isAdult = func(p person) valdy.Rule[person] {
		return valdy.NewRule(
			fmt.Sprintf("%s must be adult. age is %d", p.name, p.age),
			func(p person) bool {
				return p.age >= 18
			},
		)
	}
	hasChildren = func(p person) valdy.Rule[person] {
		return valdy.NewRule(
			fmt.Sprintf("%s should has children, but he doesn't", p.name),
			func(p person) bool {
				return len(p.children) != 0
			},
		)
	}
	inWantedList = func(criminalList []string) valdy.CreateRuleFunc[person] {
		return func(p person) valdy.Rule[person] {
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
			)
		}
	}
)

func Example() {
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
	// Output: validation error: John should has children, but he doesn't
}
