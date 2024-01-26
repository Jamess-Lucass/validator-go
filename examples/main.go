package main

import (
	"fmt"
	"strings"

	"github.com/Jamess-Lucass/validator-go"
)

type User struct {
	Firstname string
	Lastname  string
}

func main() {
	user := User{
		Firstname: "john",
		Lastname:  "doe",
	}

	rules := validator.Rules(
		validator.RuleFor(
			user.Firstname,
			validator.Required(),
			validator.Min(4),
			Custom(),
		).WithName("firstname"),
		validator.RuleFor(
			user.Lastname,
			validator.Required(),
			validator.Min(3),
		).WithName("lastname"),
	)

	validationResult := rules.Validate()

	fmt.Printf("valid: %t with errors: %v\n", validationResult.IsValid(), validationResult.Errors)
}

func Custom() *validator.RuleConfig {
	return &validator.RuleConfig{
		MessageFunc: func(r validator.RuleConfig) string {
			return fmt.Sprintf("'%s' must equal 'jane'", r.Value)
		},
		ValidateFunc: func(r validator.RuleConfig) bool {
			val, ok := r.Value.(string)
			if !ok {
				return false
			}

			return strings.EqualFold(val, "jane")
		},
	}
}
