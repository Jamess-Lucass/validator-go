package validator_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
)

type User struct {
	Firstname string `json:"firstname"`
}

func Test_MinWithValidLength(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	res := validator.RuleFor(
		user.Firstname,
		validator.Min(4),
	)

	assert.Equal(t, true, res.IsValid)
}

func Test_MinWithInvalidLength(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	res := validator.RuleFor(
		user.Firstname,
		validator.Min(100),
	)

	assert.Equal(t, false, res.IsValid)
}

func Test_MaxWithValidLength(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	res := validator.RuleFor(
		user.Firstname,
		validator.Max(5),
	)

	assert.Equal(t, true, res.IsValid)
}

func Test_MaxWithInvalidLength(t *testing.T) {
	user := User{
		Firstname: "john",
	}

	res := validator.RuleFor(
		user.Firstname,
		validator.Max(2),
	)

	assert.Equal(t, false, res.IsValid)
}
