package validator_test

import (
	"testing"

	"github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Firstname string
	Lastname  string
}

func Test_MinWithValidLength_ReturnsTrue(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	rule := validator.RuleFor(
		user.Firstname,
		validator.Min(4),
	)

	res := rule.Validate()

	assert.Equal(t, true, res.IsValid())
}

func Test_MinWithInvalidLength_ReturnsFalse(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	rule := validator.RuleFor(
		user.Firstname,
		validator.Min(100),
	)

	assert.Equal(t, false, rule.Validate().IsValid())
}

func Test_MultipleMinWithInValidLength_ReturnsFalse(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	rules := validator.Rules(
		validator.RuleFor(
			user.Firstname,
			validator.Min(4),
		),
		validator.RuleFor(
			user.Lastname,
			validator.Min(4),
		),
	)

	assert.Equal(t, false, rules.Validate().IsValid())
}

func Test_MaxWithValidLength_ReturnsTrue(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	rule := validator.RuleFor(
		user.Firstname,
		validator.Max(5),
	)

	assert.Equal(t, true, rule.Validate().IsValid())
}

func Test_MaxWithInvalidLength_ReturnsFalse(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	rule := validator.RuleFor(
		user.Firstname,
		validator.Max(2),
	)

	assert.Equal(t, false, rule.Validate().IsValid())
}

func Test_RequiredWithNoValue_ReturnsFalse(t *testing.T) {
	user := User{}

	rule := validator.RuleFor(
		user.Firstname,
		validator.Required(),
		validator.Min(2),
	)

	assert.Equal(t, false, rule.Validate().IsValid())
}
