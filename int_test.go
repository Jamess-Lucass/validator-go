package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type intTestStruct struct {
	Age      int  `json:"age"`
	Priority *int `json:"priority"`
}

func TestInt_NotEmpty(t *testing.T) {
	s := intTestStruct{Age: 0}
	v, _ := validator.New(&s)
	v.Int(&s.Age).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 5
	v, _ = validator.New(&s)
	v.Int(&s.Age).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Gte(t *testing.T) {
	s := intTestStruct{Age: 17}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Gte(18)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 18
	v, _ = validator.New(&s)
	v.Int(&s.Age).Gte(18)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Lte(t *testing.T) {
	s := intTestStruct{Age: 121}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Lte(120)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 120
	v, _ = validator.New(&s)
	v.Int(&s.Age).Lte(120)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Gt(t *testing.T) {
	s := intTestStruct{Age: 0}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Gt(0)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 1
	v, _ = validator.New(&s)
	v.Int(&s.Age).Gt(0)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Lt(t *testing.T) {
	s := intTestStruct{Age: 100}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Lt(100)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 99
	v, _ = validator.New(&s)
	v.Int(&s.Age).Lt(100)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Positive(t *testing.T) {
	s := intTestStruct{Age: -1}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Positive()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 1
	v, _ = validator.New(&s)
	v.Int(&s.Age).Positive()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Negative(t *testing.T) {
	s := intTestStruct{Age: 1}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Negative()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = -1
	v, _ = validator.New(&s)
	v.Int(&s.Age).Negative()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_MultipleOf(t *testing.T) {
	s := intTestStruct{Age: 7}
	v, _ := validator.New(&s)
	v.Int(&s.Age).MultipleOf(3)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 9
	v, _ = validator.New(&s)
	v.Int(&s.Age).MultipleOf(3)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Nonnegative(t *testing.T) {
	s := intTestStruct{Age: -1}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Nonnegative()
	assert.False(t, v.Validate().IsValid())

	s.Age = 0
	v, _ = validator.New(&s)
	v.Int(&s.Age).Nonnegative()
	assert.True(t, v.Validate().IsValid())
}

func TestInt_Nonpositive(t *testing.T) {
	s := intTestStruct{Age: 1}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Nonpositive()
	assert.False(t, v.Validate().IsValid())

	s.Age = 0
	v, _ = validator.New(&s)
	v.Int(&s.Age).Nonpositive()
	assert.True(t, v.Validate().IsValid())
}

func TestInt_WithName(t *testing.T) {
	s := intTestStruct{Priority: nil}
	v, _ := validator.New(&s)
	v.Int(s.Priority).NotNil().WithName("priority")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "priority", result.Errors[0].Field)
}

func TestInt_NilSkip(t *testing.T) {
	s := intTestStruct{Priority: nil}
	v, _ := validator.New(&s)
	v.Int(s.Priority).Gte(1)
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_NotNil(t *testing.T) {
	s := intTestStruct{Priority: nil}
	v, _ := validator.New(&s)
	v.Int(s.Priority).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
}

func TestInt_Must(t *testing.T) {
	s := intTestStruct{Age: 13}
	v, _ := validator.New(&s)
	v.Int(&s.Age).Must(func(val int) bool {
		return val%2 == 0
	}).WithMessage("must be even")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be even", result.Errors[0].Message)
}

func TestInt_Standalone(t *testing.T) {
	rule := validator.Int().Gte(18)

	errs := rule.Validate(20)
	assert.Empty(t, errs)

	errs = rule.Validate(10)
	assert.Len(t, errs, 1)

	// JSON numbers decode as float64
	errs = rule.Validate(float64(20))
	assert.Empty(t, errs)

	errs = rule.Validate("not a number")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)
}

func TestInt_Coercion_AllNumericKinds(t *testing.T) {
	rule := validator.Int().Gte(18)

	// Signed, unsigned, and integral floats of every width are accepted.
	assert.Empty(t, rule.Validate(int8(20)))
	assert.Empty(t, rule.Validate(int16(20)))
	assert.Empty(t, rule.Validate(int32(20)))
	assert.Empty(t, rule.Validate(int64(20)))
	assert.Empty(t, rule.Validate(uint(20)))
	assert.Empty(t, rule.Validate(uint8(20)))
	assert.Empty(t, rule.Validate(uint64(20)))
	assert.Empty(t, rule.Validate(float32(20)))
	assert.Empty(t, rule.Validate(float64(20)))

	// A float with a fractional part is not an integer.
	errs := rule.Validate(float64(20.5))
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)
}
