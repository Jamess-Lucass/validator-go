package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

type boolTestStruct struct {
	Active   bool  `json:"active"`
	Accepted *bool `json:"accepted"`
}

func TestBool_NotNil_WithName(t *testing.T) {
	s := boolTestStruct{Accepted: nil}
	v, _ := validator.New(&s)
	validator.Bool(v, s.Accepted).NotNil().WithName("accepted")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "accepted", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestBool_NotEmpty(t *testing.T) {
	s := boolTestStruct{Active: false}

	v, _ := validator.New(&s)
	validator.Bool(v, &s.Active).NotEmpty()

	result := v.Validate()

	assert.False(t, result.IsValid())

	s.Active = true

	v, _ = validator.New(&s)
	validator.Bool(v, &s.Active).NotEmpty()

	result = v.Validate()

	assert.True(t, result.IsValid())
}

func TestBool_Must(t *testing.T) {
	s := boolTestStruct{Active: false}

	v, _ := validator.New(&s)

	validator.Bool(v, &s.Active).Must(func(val bool) bool {
		return val == true
	}).WithMessage("must accept terms")

	result := v.Validate()

	assert.False(t, result.IsValid())
	assert.Equal(t, "must accept terms", result.Errors[0].Message)
}

func TestBool_NamedBoolType(t *testing.T) {
	type flag bool
	assert.Empty(t, rule.Bool().NotEmpty().Validate(flag(true)))
	assert.Len(t, rule.Bool().NotEmpty().Validate(flag(false)), 1)
}

func TestBool_NamedBoolType_StructMode(t *testing.T) {
	type flag bool
	type settings struct {
		Beta flag `json:"beta"`
	}

	s := settings{Beta: false}
	v, _ := validator.New(&s)
	validator.Bool(v, &s.Beta).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "beta", result.Errors[0].Field)

	s.Beta = true
	v, _ = validator.New(&s)
	validator.Bool(v, &s.Beta).NotEmpty()
	assert.True(t, v.Validate().IsValid())
}

func TestBool_Standalone(t *testing.T) {
	r := rule.Bool().NotEmpty()

	errs := r.Validate(true)
	assert.Empty(t, errs)

	errs = r.Validate(false)
	assert.Len(t, errs, 1)

	errs = r.Validate("not a bool")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a boolean", errs[0].Message)
}
