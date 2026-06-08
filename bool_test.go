package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type boolTestStruct struct {
	Active   bool  `json:"active"`
	Accepted *bool `json:"accepted"`
}

func TestBool_NotNil_WithName(t *testing.T) {
	s := boolTestStruct{Accepted: nil}
	v, _ := validator.New(&s)
	v.Bool(s.Accepted).NotNil().WithName("accepted")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "accepted", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestBool_NotEmpty(t *testing.T) {
	s := boolTestStruct{Active: false}

	v, _ := validator.New(&s)
	v.Bool(&s.Active).NotEmpty()

	result := v.Validate()

	assert.False(t, result.IsValid())

	s.Active = true

	v, _ = validator.New(&s)
	v.Bool(&s.Active).NotEmpty()

	result = v.Validate()

	assert.True(t, result.IsValid())
}

func TestBool_Must(t *testing.T) {
	s := boolTestStruct{Active: false}

	v, _ := validator.New(&s)

	v.Bool(&s.Active).Must(func(val bool) bool {
		return val == true
	}).WithMessage("must accept terms")

	result := v.Validate()

	assert.False(t, result.IsValid())
	assert.Equal(t, "must accept terms", result.Errors[0].Message)
}

func TestBool_Standalone(t *testing.T) {
	rule := validator.Bool().NotEmpty()

	errs := rule.Validate(true)
	assert.Empty(t, errs)

	errs = rule.Validate(false)
	assert.Len(t, errs, 1)

	errs = rule.Validate("not a bool")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a boolean", errs[0].Message)
}
